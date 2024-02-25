package errors

type Reporter = func(line int, message string)
