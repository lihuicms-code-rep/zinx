package znet

import (
	"fmt"
	"io"
	"net"
	"testing"
)

//负责测试datapack的单元测试
func TestDataPack(t *testing.T) {
	//1.创建socket句柄
	listener, err := net.Listen("tcp","0.0.0.0:7777" )

	if err != nil {
		fmt.Println("server listen err:", err)
		return
	}

	//2.承载业务
	go func() {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("accept error:", err)
		}

		//读取客户端数据
		go func(conn net.Conn) {
			//拆包过程
			for {
				//1.第一次,从conn读取出head
				dp := NewDataPack()
				headData := make([]byte, dp.GetHeadLen())
				_, err := io.ReadFull(conn, headData)
				if err != nil {
					fmt.Println("read head error", err)
				}

				msgHead, err := dp.Unpack(headData)
				if err != nil {
					return
				}

				if msgHead.GetMsgLen() > 0 {
					//2.第二次,从conn继续读出data
					msg := msgHead.(*Message)
					msg.Data = make([]byte, msg.GetMsgLen())
					_, err := io.ReadFull(conn, msg.Data)
					if err != nil {
						fmt.Println("unpack data error", err)
						return
					}


					//完整消息读取完毕
					fmt.Println("Receive MsgId:", msg.Id, " dataLen:", msg.DataLen, " data=", msg.Data)
				}
			}
		}(conn)
	}()


	//3.模拟客户端
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client dial err", err)
		return
	}

	dp := NewDataPack()

	//模拟粘包过程
	msg1 := &Message{
		Id:1,
		DataLen:5,
		Data:[]byte{'H', 'E', 'L', 'L', 'O'},
	}

	sendData1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("pack msg1 error", err)
		return
	}


	msg2 := &Message {
		Id:2,
		DataLen:4,
		Data:[]byte{'z', 'i', 'n', 'x'},
	}

	sendData2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println("pack msg2 error", err)
		return
	}

	sendData1 = append(sendData1, sendData2...)
	conn.Write(sendData1)

	//客户端阻塞
	select{}
}


