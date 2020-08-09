package arest

import (
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

type ArestTestSuite struct {
	suite.Suite
	client Arest
}

func (s *ArestTestSuite) SetupSuite() {
	// Init logger
	logrus.SetFormatter(new(prefixed.TextFormatter))
	logrus.SetLevel(logrus.DebugLevel)
}
