package service

import (
	"math/rand"
	"time"
)

var defaultAvatars = []string{
	"https://001-1256564327.cos.ap-beijing.myqcloud.com/2022/05/30/default-boy.png",
	"https://001-1256564327.cos.ap-beijing.myqcloud.com/2022/05/30/default-girl.png",
}

func (s *Service) GetRandomAvatar() string {
	rand.Seed(time.Now().UnixMicro())
	return defaultAvatars[rand.Intn(len(defaultAvatars))]
}
