package main

import (
	"fmt"
	
	"github.com/wanglu119/me-deps/viper"
)

func main() {
	vp := viper.GetViper()
	vp.Set("test", "value")
	
	fmt.Println(vp.GetString("test"))
}
