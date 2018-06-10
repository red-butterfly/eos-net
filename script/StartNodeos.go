package script

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"../util"
)

func StartEOSNode() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))  //返回绝对路径  filepath.Dir(os.Args[0])去除最后一个元素的路径
	if err != nil {
		fmt.Println(err)
	}

	nodeargs := []string{
		strings.Join([]string{"/tmp", strings.ToUpper(util.ActInfo.NodeProducer[0])}, "/") ,
		util.ActInfo.NodeProducer[0],
		util.ActInfo.NodeProducer[2],
		util.ActInfo.NodeProducer[1],
		dir,
	}

	cmd := exec.Command("./script/nodeos_genesis.sh", nodeargs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err = cmd.Run()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("---- END ----\n")
}
