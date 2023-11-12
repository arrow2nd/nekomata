package app

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/arrow2nd/nekomata/api/sharedapi"
	"github.com/arrow2nd/nekomata/cli"
	"github.com/skanehira/clipboard-image/v2"
	"github.com/spf13/pflag"
	"golang.org/x/sync/errgroup"
	"golang.org/x/term"
)

func (a *App) newPostCmd() *cli.Command {
	longHelp := `Submit a post.
When the post content is abbreviated, the external editor set to $EDITOR will be launched.`

	example := `post にゃーん --image cat.png dog.png
  echo "にゃーん" | nekomata post`

	return &cli.Command{
		Name:      "post",
		Short:     "Post",
		Long:      longHelp,
		UsageArgs: "[text]",
		Example:   example,
		SetFlag: func(f *pflag.FlagSet) {
			// 投稿に関する追加情報
			f.StringP("reply", "r", "", "reply destination post ID")
			f.StringP("visibility", "v", "public", "post visibility range")
			f.BoolP("nsfw", "n", false, "set the NSFW flag")

			// コマンドの挙動関係
			f.StringP("editor", "e", os.Getenv("EDITOR"), "external editor for editing")

			// 添付メディア関係
			f.StringSliceP("image", "i", nil, "path of the image(s) to attach (multiple can be specified, separated by commas)")
			f.BoolP("clipboard", "c", false, "attach the image from the clipboard (this is ignored if --image is specified)")
		},
		Run: a.execPostCmd,
	}
}

func (a *App) execPostCmd(c *cli.Command, f *pflag.FlagSet) error {
	text := ""

	if f.NArg() == 0 && !term.IsTerminal(int(syscall.Stdin)) {
		// 標準入力から受け取る
		stdin, _ := io.ReadAll(os.Stdin)
		text = string(stdin)
	} else {
		// 引数をすべてスペースで連結
		text = strings.Join(f.Args(), " ")
	}

	// 外部エディタを起動
	if text == "" {
		editor, _ := f.GetString("editor")

		t, err := a.editPostWithExternalEditor(editor)
		if err != nil {
			return err
		}

		text = t
	}

	submitPost(f, text)
	return nil
}

// editPostWithExternalEditor : 投稿を外部エディタで編集する
func (a *App) editPostWithExternalEditor(editor string) (string, error) {
	tmpFilePath := path.Join(os.TempDir(), ".nekomata_tmp")
	if _, err := os.Create(tmpFilePath); err != nil {
		return "", err
	}

	if err := a.openExternalEditor(editor, tmpFilePath); err != nil {
		return "", err
	}

	bytes, err := os.ReadFile(tmpFilePath)
	if err != nil {
		return "", err
	}

	if err := os.Remove(tmpFilePath); err != nil {
		return "", err
	}

	return string(bytes), nil
}

// submitPost : 投稿を送信
func submitPost(f *pflag.FlagSet, t string) {
	images, _ := f.GetStringSlice("image")
	text := strings.TrimSpace(t)

	// 本文と画像が無い場合はキャンセル
	if text == "" && len(images) == 0 {
		return
	}

	nsfw, _ := f.GetBool("nsfw")
	visibility, _ := f.GetString("visibility")
	postOpts := &sharedapi.CreatePostOpts{
		Text:       text,
		Visibility: visibility,
		Sensitive:  nsfw,
	}

	replyID, _ := f.GetString("reply")
	existClipboardImage, _ := f.GetBool("clipboard")

	doPost := func() {
		if existImages := len(images) > 0; existImages || existClipboardImage {
			var err error

			if existImages {
				postOpts.MediaIDs, err = uploadImages(images)
			} else {
				postOpts.MediaIDs, err = uploadImageFromClipboard()
			}

			if err != nil {
				global.SetErrorStatus("Upload media", err.Error())
				return
			}
		}

		res, err := global.client.CreatePost(postOpts)
		if err != nil {
			global.SetErrorStatus("Post", err.Error())
		}

		statusLabel := "Posted"
		if count := len(postOpts.MediaIDs); count > 0 {
			statusLabel += fmt.Sprintf(" / with %d attached images", count)
		}

		global.SetStatus(statusLabel, res.Text)
	}

	// CLIなら実行
	if global.isCLI {
		doPost()
		return
	}

	// 実行しようとしている操作名
	operation := "post"
	if replyID != "" {
		operation = "reply"
	}

	// TODO: 独自モーダルを出す
	global.ReqestPopupModal(&ModalOpts{
		title:  fmt.Sprintf("Do you want to post a [%s]%s[-:-:-]?", global.conf.Style.App.EmphasisText, operation),
		text:   text,
		onDone: doPost,
	})
}

// uploadImageFromClipboard : クリップボードの画像をアップロード
func uploadImageFromClipboard() ([]string, error) {
	r, err := clipboard.Read()
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, r); err != nil {
		return nil, err
	}

	id, err := global.client.UploadMedia("image", buf)
	if err != nil {
		return nil, fmt.Errorf("upload failed: %w", err)
	}

	return []string{id}, nil
}

// uploadImages : 複数の画像をアップロード
func uploadImages(paths []string) ([]string, error) {
	imageCounts := len(paths)

	// 画像の枚数チェック
	if imageCounts > 4 {
		return nil, errors.New("you can attach up to 4 images")
	}

	eg, ctx := errgroup.WithContext(context.Background())
	ch := make(chan string, imageCounts)

	for _, imagePath := range paths {
		// 拡張子のチェック
		ext := strings.ToLower(path.Ext(imagePath))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" {
			return nil, fmt.Errorf("unsupported extensions (%s)", imagePath)
		}

		imagePath := imagePath

		eg.Go(func() error {
			select {
			case <-ctx.Done():
				return nil
			default:
				f, err := os.Open(imagePath)
				if err != nil {
					return fmt.Errorf("failed to load file (%s)", imagePath)
				}

				defer f.Close()

				id, err := global.client.UploadMedia(filepath.Base(imagePath), f)
				if err != nil {
					return fmt.Errorf("upload failed (%s): %w", imagePath, err)
				}

				ch <- id
				return nil
			}
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	close(ch)

	mediaIDs := []string{}
	for id := range ch {
		mediaIDs = append(mediaIDs, id)
	}

	return mediaIDs, nil
}
