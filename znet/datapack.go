package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"zinx/utils"
	"zinx/ziface"
)

type DataPack struct{}

//实例对象
func NewDataPack() *DataPack {
	return &DataPack{}
}

func (dp *DataPack) GetHeadLen() uint32 {
	return 4 + 4 //id所占字节+dataLen所占字节
}

//封包
func (dp *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	dataBuf := bytes.NewBuffer([]byte{}) //存放byte字节的缓冲

	//注意所写字段的顺序

	//将dataLen写入dataBuf
	if err := binary.Write(dataBuf, binary.LittleEndian, msg.GetMsgLen()); err != nil {
		return nil, err
	}

	//将msgId写入dataBuf
	if err := binary.Write(dataBuf, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}

	//将msgData写入dataBuf
	if err := binary.Write(dataBuf, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}

	return dataBuf.Bytes(), nil
}

//拆包:只需要将头部分读出来,之后就可以根据head提供的信息读出data
func (dp *DataPack) Unpack(data []byte) (ziface.IMessage, error) {
	dataBuf := bytes.NewReader(data) //创建一个二进制reader
	msg := &Message{}
	//先读dataLen到msg
	if err := binary.Read(dataBuf, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}

	//再读msgId到msg
	if err := binary.Read(dataBuf, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}

	//判断dataLen是否超过最大包长度
	if utils.GlobalObject.MaxPackageSize > 0 && msg.DataLen > utils.GlobalObject.MaxPackageSize {
		return nil, errors.New("too large msg data")
	}

	return msg, nil
}
