package utils

import (
	"fmt"

	"google.golang.org/protobuf/reflect/protoreflect"
)

// GetProtoMethods 获取 proto 文件中的所有方法
func GetProtoMethods(fileDescriptor protoreflect.FileDescriptor) (list []string) {
	services := fileDescriptor.Services()
	for i := 0; i < services.Len(); i++ {
		service := services.Get(i)
		methods := service.Methods()
		for j := 0; j < methods.Len(); j++ {
			method := methods.Get(j)

			item := fmt.Sprintf("%s/%s", service.FullName(), method.Name())
			list = append(list, item)
		}
	}
	return
}
