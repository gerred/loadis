package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/siddontang/ledisdb/ledis"
)

type Command struct {
	Text string

	query string
}

func (cmd *Command) Execute(ledisInfo *LedisInfo) error {
	text := strings.ToLower(cmd.Text)
	if text == "keys" {
		return cmd.runKeys(ledisInfo)
	}

	cmd.parseQuery()

	if strings.HasPrefix(text, "hgetall") {
		return cmd.runHGetAll(ledisInfo)
	}
	if strings.HasPrefix(text, "hget") {
		return cmd.runHGet(ledisInfo)
	}
	if strings.HasPrefix(text, "smembers") {
		return cmd.runSmembers(ledisInfo)
	}
	if strings.HasPrefix(text, "get") {
		return cmd.runGet(ledisInfo)
	}

	return fmt.Errorf("unknown command: %s\n", cmd.Text)
}

func (cmd *Command) parseQuery() {
	parts := strings.SplitN(cmd.Text, " ", 2)
	if len(parts) == 2 {
		cmd.query = parts[1]
	}
}

func (cmd *Command) runKeys(ledisInfo *LedisInfo) error {
	keyTypes := []ledis.DataType{
		ledis.KV, ledis.SET, ledis.LIST, ledis.HASH, ledis.ZSET,
	}
	for _, keyType := range keyTypes {
		keys, err := ledisInfo.GetKeyList(keyType)
		if err != nil {
			return err
		}

		fmt.Printf("======keys of type %v\n", keyType)
		for _, key := range keys {
			fmt.Printf("%s\n", key)
		}
	}

	return nil
}

func (cmd *Command) runHGetAll(ledisInfo *LedisInfo) error {
	vals, err := ledisInfo.Db.HGetAll([]byte(cmd.query))
	if err != nil {
		return err
	}
	for _, val := range vals {
		fmt.Printf("=======field:\n")
		fmt.Printf("%s:%s\n", val.Field, val.Value)
	}
	return nil
}

func (cmd *Command) runHGet(ledisInfo *LedisInfo) error {
	return errors.New("not implemented")
}

func (cmd *Command) runSmembers(ledisInfo *LedisInfo) error {
	vals, err := ledisInfo.Db.SMembers([]byte(cmd.query))
	if err != nil {
		return err
	}
	for _, val := range vals {
		fmt.Printf("%s\n", val)
	}
	return nil
}

func (cmd *Command) runGet(ledisInfo *LedisInfo) error {
	val, err := ledisInfo.Db.Get([]byte(cmd.query))
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", val)
	return nil
}
