package stringx

func InterfaceTrans(raw interface{}, v interface{}) error {
	str, err := Marshal(raw)
	if err != nil {
		return err
	}
	err = Unmarshal(str, &v)
	return err
}
