package view

import (
	"context"
	"fmt"
	"gototp/pkg/timeutil"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
)

func (v *View) withtime() error {
	var input string

	form := huh.NewForm(
		huh.NewGroup(
			Input(
				"time ", timeutil.NowS(),
				strValidator,
				&input,
				timeutil.NowS(),
			),
		).WithTheme(base16),
	)
	if err := form.Run(); err != nil {
		return err
	}

	t, err := timeutil.Parse(input)
	if err != nil {
		return err
	}
	delta := time.Since(t)

	ctx, cancel := context.WithCancel(context.Background())

	id, err := v.SelectCode(fmt.Sprintf("🔑 generate code with time %s", timeutil.Format(t)))
	if err != nil {
		cancel()
		return err
	}
	key := v.gototp.Data.Keystore.Keys[id]
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				otpcode, err := key.GenCode(time.Now().Add(-delta))
				if err != nil {
					ErrorFunc(err)
					cancel()
					return
				}
				fmt.Printf("🔑 %s - %s\r\n", key.Name, otpcode)
				time.Sleep(time.Second * time.Duration(key.Period))
			}
		}
	}()
	spinner.New().Title(" updating").Context(ctx).Type(spinner.Monkey).Run()
	cancel()
	return nil
}
