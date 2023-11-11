package znet

import (
	"errors"
	"fmt"
	"sync"

	"zinx.mod/ziface"
)

// 连接管理模块
type ConnManager struct {
	Connections map[uint32]ziface.IConnection //连接集合
	ConnLock    sync.RWMutex                  //保护连接集合的读写锁
}

func NewConnManager() ziface.IConnManager {
	Connmrg := &ConnManager{
		Connections: make(map[uint32]ziface.IConnection),
	}
	return Connmrg
}

// 添加连接
func (cm *ConnManager) Add(conn ziface.IConnection) {
	cm.ConnLock.Lock()
	defer cm.ConnLock.Unlock()

	cm.Connections[conn.GetConnId()] = conn

	fmt.Println("ConnId is ", conn.GetConnId(), "  Add to Connections Successful,Connection len is ", len(cm.Connections))
}

// 删除连接
func (cm *ConnManager) Remote(conn ziface.IConnection) {
	cm.ConnLock.Lock()
	defer cm.ConnLock.Unlock()

	delete(cm.Connections, conn.GetConnId())
	fmt.Println("ConnId is ", conn.GetConnId(), "  Remote to Connections Successful,Connection len is ", len(cm.Connections))
}

// 根据connID获取连接
func (cm *ConnManager) Get(connID uint32) (ziface.IConnection, error) {
	cm.ConnLock.RLock()
	defer cm.ConnLock.RUnlock()

	if conn, ok := cm.Connections[connID]; ok {
		return conn, nil
	} else {
		return nil, errors.New("ConnID Not Exist")
	}
}

// 得到当前连接的总条数
func (cm *ConnManager) Len() int {
	return len(cm.Connections)
}

// 清除并终止当前所有连接
func (cm *ConnManager) ClearConn() {
	cm.ConnLock.Lock()
	defer cm.ConnLock.Unlock()

	for connID, conn := range cm.Connections {
		//停止
		conn.Stop()
		//删除
		delete(cm.Connections, connID)
	}

	fmt.Println("Clear Conn Successful , map len is ", len(cm.Connections))
}
