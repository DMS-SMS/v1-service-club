package handler

import (
	"club/db"
	authproto "club/proto/golang/auth"
	"club/tool/consul"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/opentracing/opentracing-go"
)

type _default struct {
	None
	accessManage db.AccessorManage
	tracer       opentracing.Tracer
	awsSession   *session.Session
	consulAgent  consul.Agent
	authStudent  authproto.AuthStudentService
}

type FieldSetter func(*_default)

func Default(setters ...FieldSetter) *_default {
	return newDefault(setters...)
}

func newDefault(setters ...FieldSetter) (h *_default) {
	h = new(_default)
	for _, setter := range setters {
		setter(h)
	}
	return
}

func AccessManager(am db.AccessorManage) FieldSetter {
	return func(h *_default) {
		h.accessManage = am
	}
}

func Tracer(t opentracing.Tracer) FieldSetter {
	return func(h *_default) {
		h.tracer = t
	}
}

func AWSSession(s *session.Session) FieldSetter {
	return func(h *_default) {
		h.awsSession = s
	}
}

func ConsulAgent(ca consul.Agent) FieldSetter {
	return func(h *_default) {
		h.consulAgent = ca
	}
}

func AuthStudent(as authproto.AuthStudentService) FieldSetter {
	return func(h *_default) {
		h.authStudent = as
	}
}
