package help

func Text(command string) string {
	bindata, _ := Asset(command + `.fmt`)
	return string(bindata) + "\n"
}
