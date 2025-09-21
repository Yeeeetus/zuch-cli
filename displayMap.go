package main

var testMap = []string{
	/*
	 0.1.2.3.4.5.6.7.8.9*/
	"-.+.-.-.+.-.-.-.-.-", //0
	" .|. . .+.+. . . . ", //1
	" . . . . .+.+. . . ", //2
	" .|. . . . .|. . . ", //3
	" .+.-.-.+.-.+. . . ", //4
	" . . . .|. .|. . . ", //5
	" . . . .|. .|. . . ", //6
	" . .-.-.+.-.+.-.-. ", //7
	" . . . . . .|. . . ", //8
	" . . . . . .|. . . ", //9
}

func convertMapToString(model model) string {
	result := ""
	for k, v := range testMap {
		if k == 0 {

			result += v
		} else {

			result += "\n" + v
		}
	}
	return result
}
