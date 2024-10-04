package view

import (
	"context"
	"fmt"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
)

func (v *View) withtime() error {
	var input string

	form := huh.NewForm(
		huh.NewGroup(
			Input(
				"time ", time.Now().Format(timeFormat),
				strValidator,
				&input,
				fmt.Sprint(time.Now().Format(timeFormat)),
			),
		).WithTheme(base16),
	)
	if err := form.Run(); err != nil {
		return err
	}

	with, err := time.Parse(timeFormat, input)
	if err != nil {
		return fmt.Errorf("error parsing time: %w", err)
	}

	delta := time.Since(with)
	fmt.Println(delta.Abs())

	ctx, cancel := context.WithCancel(context.Background())

	id, err := v.SelectCode(fmt.Sprintf("🔑 generate code with time %s", with.Format(timeFormat)))
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
				fmt.Println(time.Now().Add(-delta).Format(timeFormat))

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
