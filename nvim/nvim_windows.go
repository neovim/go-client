package nvim

import "syscall"

func init() {
	embedProcAttr = &syscall.SysProcAttr{HideWindow: true}
}
