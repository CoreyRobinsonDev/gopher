package commands

func Unwrap[T any](val T, err error) T {
	if err != nil { panic(err) }

	return val
}

func UnwrapOr[T any](val T, err error) func(T) T {
	if err != nil {
		return func(d T) T {
			return d
		}
	} else {
		return func(_ T) T {
			return val
		}
	}
}

func Expect(err error) {
	if err != nil { panic(err) }
}

