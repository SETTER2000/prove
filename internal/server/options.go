package server

import (
	"fmt"
	"github.com/SETTER2000/prove/config"
	"github.com/SETTER2000/prove/internal/controller/grpc"
	"github.com/SETTER2000/prove/internal/controller/grpc/handler"
	"github.com/SETTER2000/prove/scripts"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/acme/autocert"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

// Option -.
type Option func(*Server)

// Port -.
//func Port(port string) Option {
//	return func(s *Server) {
//		s.server.Addr = net.JoinHostPort("", port)
//	}
//}

//// Host -.
//func Host(host string) Option {
//	return func(s *Server) {
//		domain := strings.Split(host, ":")
//		s.server.Addr = net.JoinHostPort(domain[0], domain[1])
//	}
//}

// Host -.
func Host() Option {
	host := config.GetConfig().ServerAddress
	log.Printf("Server is listening at: %s", host)
	return func(s *Server) {
		domain := strings.Split(host, ":")
		s.server.Addr = net.JoinHostPort(domain[0], domain[1])
	}
}

//func Host(host string) Option {
//	return func(s *Server) {
//		domain := strings.Split(host, ":")
//		s.server.Addr = net.JoinHostPort(domain[0], domain[1])
//	}
//}

// ReadTimeout -.
func ReadTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.server.ReadTimeout = timeout
	}
}

// WriteTimeout -.
func WriteTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.server.WriteTimeout = timeout
	}
}

// ShutdownTimeout -.
func ShutdownTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.shutdownTimeout = timeout
	}
}

// PortGRPC -.
func PortGRPC() Option {
	return func(s *Server) {
		s.grpcPort = config.GetConfig().Port
	}
}

// EnableGRPC - включить поддержку gRPC.
func EnableGRPC(h *handler.IProveServer) Option {
	var logger = logrus.New()

	grpcSrv := grpc.NewServer(grpc.Deps{
		Logger:  logger,
		Handler: h,
	})

	return func(s *Server) {
		s.srv = grpcSrv
	}
}

// EnableHTTPS - опция подключает возможность использования SSL/TLS на сервере.
func EnableHTTPS() Option {
	cfg := config.GetConfig().HTTP
	fmt.Printf("TLS cert. ServerDomain: %s\n", config.GetConfig().ServerDomain)
	fmt.Printf("CONNECT DB DSN: %s\n", config.GetConfig().ConnectDB)
	// конструируем менеджер TLS-сертификатов
	manager := &autocert.Manager{
		// директория для хранения сертификатов
		Cache: autocert.DirCache(config.GetConfig().CertsDir),
		// функция, принимающая Terms of Service издателя сертификатов
		Prompt: autocert.AcceptTOS,
		// перечень доменов, для которых будут поддерживаться сертификаты
		HostPolicy: autocert.HostWhitelist(config.GetConfig().ServerDomain, fmt.Sprintf("www.%s", config.GetConfig().ServerDomain)),
	}
	tlsConfig := manager.TLSConfig()
	tlsConfig.GetCertificate = scripts.GetSelfSignedOrLetsEncryptCert(manager)
	return func(s *Server) {
		flag := config.GetConfig().EnableHTTPS
		_, ok := os.LookupEnv("ENABLE_HTTPS")

		if !ok {
			if flag {
				s.isHTTPS = flag
				s.server.Addr = ":443"
				s.server.TLSConfig = tlsConfig
				s.certFile = fmt.Sprintf("%s/%s.crt", cfg.CertsDir, config.GetConfig().ServerDomain)
				s.keyFile = fmt.Sprintf("%s/%s.key", cfg.CertsDir, config.GetConfig().ServerDomain)
				log.Printf("enabled HTTPS: %v\n", flag)
			} else {
				s.isHTTPS = flag
				s.certFile = ""
				s.keyFile = ""
				log.Printf("disabled HTTPS: %v\n", flag)
			}
		} else {
			s.isHTTPS = flag
			log.Printf("connect HTTPS: %v\n", flag)
		}
	}
}
