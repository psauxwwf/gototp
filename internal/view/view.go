package view

import (
	"errors"
	"fmt"
	"gototp/internal/gototp"
	"gototp/pkg/fs"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/huh"

	"github.com/pquerna/otp"
)

var (
	strValidator = func(in string) error {
		if in == "" {
			return fmt.Errorf("⚠️  value is empty")
		}
		return nil
	}
	numValidator = func(in string) error {
		if _, err := strconv.Atoi(in); err != nil {
			return fmt.Errorf("⚠️  value must be digit")
		}
		return nil
	}
)

type View struct {
	gototp *gototp.Gototp
}

func New() *View {
	_gototp, err := unlock(readPassphrase)
	if err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			os.Exit(0)
		}
		ErrorFunc(err)
		os.Exit(1)
	}

	v := View{
		gototp: _gototp,
	}
	return &v
}

const (
	_show     = "📒 show"
	_code     = "🔑 code"
	_add      = "🆕 add"
	_delete   = "❌ delete"
	_export   = "💾 export"
	_import   = "📥 import"
	_change   = "🔄 change"
	_withtime = "⏰ code for time"
)

var mainopts = []huh.Option[string]{
	huh.NewOption[string](_show, _show),
	huh.NewOption[string](_code, _code),
	huh.NewOption[string](_add, _add),
	huh.NewOption[string](_change, _change),
	huh.NewOption[string](_delete, _delete),
	huh.NewOption[string](_export, _export),
	huh.NewOption[string](_import, _import),
	huh.NewOption[string](_withtime, _withtime),
}

var (
	algorithmops = huh.NewOptions[string](
		fmt.Sprint(otp.AlgorithmSHA1),
		fmt.Sprint(otp.AlgorithmSHA256),
		fmt.Sprint(otp.AlgorithmSHA512),
		fmt.Sprint(otp.AlgorithmMD5),
	)
	digitopts = huh.NewOptions[string](
		fmt.Sprint(otp.DigitsSix),
		fmt.Sprint(otp.DigitsEight),
	)
)

func (v *View) SelectCode(title string) (int, error) {
	var (
		strid string
		opts  []huh.Option[string]
	)
	if len(v.gototp.Data.Keystore.Keys) == 0 {
		return -1, fmt.Errorf("keys not created")
	}
	for id, f := range v.gototp.Data.Keystore.Keys {
		opts = append(opts, huh.NewOption[string](f.Name, fmt.Sprint(id)))
	}
	form := SelectForm(title, opts, &strid)
	if err := form.Run(); err != nil {
		return -1, err
	}
	id, err := strconv.Atoi(strid)
	if err != nil {
		return -1, fmt.Errorf("failed to get id")
	}
	return id, nil
}

func (v *View) SelectLocalFile(title string, suffix ...string) (string, error) {
	var (
		file   string
		opts   []huh.Option[string]
		sorted []string
	)
	files, err := fs.FilesInCurrentDir()
	if err != nil {
		return "", err
	}
	if len(suffix) == 1 {
		for _, file := range files {
			if strings.HasSuffix(file, suffix[0]) {
				sorted = append(sorted, file)
			}
		}
		files = sorted
	}
	for _, f := range files {
		opts = append(opts, huh.NewOption[string](f, f))
	}
	form := SelectForm(title, opts, &file)
	if err := form.Run(); err != nil {
		return "", err
	}
	return file, nil
}

func (v *View) Run() {
	actionHandlers := map[string]func() error{
		_show:     v.show,
		_code:     v.code,
		_withtime: v.withtime,
		_add:      v.add,
		_change:   v.change,
		_delete:   v.delete,
		_import:   v.importa,
		_export:   v.export,
	}

	for {
		var action string
		form := SelectForm("🔐 gototp ", mainopts, &action)
		if err := form.Run(); err != nil {
			if errors.Is(err, huh.ErrUserAborted) {
				return
			}
			ErrorFunc(err)
			continue
		}

		if handler, exists := actionHandlers[action]; exists {
			if err := handler(); err != nil {
				if errors.Is(err, huh.ErrUserAborted) {
					continue
				}
				ErrorFunc(err)
			}
		}
	}
}
