package main

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// TODO implement
func validConnectionParams(authMethod int, name, host, username, keyFile, password, port string) error {
	var err error = nil

	if err = isIP(host); err != nil {
		return err
	}
	if err = isValidAuthMethod(authMethod, keyFile, password); err != nil {
		return err
	}
	if err = isPort(port); err != nil {
		return err
	}
	if authMethod == 1 {
		if err = isValidFilePath(keyFile); err != nil {
			return err
		}
	}

	return nil
}

func isIP(ip string) error {
	octets := strings.Split(ip, ".")
	if len(octets) != 4 {
		return errors.New("not valid ip address")
	}

	for _, octet := range octets {
		if oc, err := strconv.Atoi(octet); err != nil || oc > 255 || oc < 0 {
			return errors.New("not valid ip address")
		}
	}
	return nil
}

func isValidFilePath(path string) error {
	if !filepath.IsAbs(path) {
		return errors.New("file path is not absolute")
	}

	if _, err := os.Stat(path); err != nil {
		return errors.New("file path may not exist")
	}
	return nil
}

func isValidAuthMethod(method int, key, passw string) error {
	if method == 0 {
		if passw == "" {
			return errors.New("password cannot be empty")
		}
	} else {
		if key == "" {
			return errors.New("public key cannot be empty")
		}
	}
	return nil
}

func isPort(s string) error {
	if _, err := strconv.Atoi(s); err != nil {
		return errors.New("port should be number")
	}
	return nil
}
