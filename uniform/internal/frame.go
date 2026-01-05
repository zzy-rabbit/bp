package internal

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"github.com/zzy-rabbit/bp/uniform/api"
	"github.com/zzy-rabbit/xtools/xerror"
	"io"
)

func (s *service) MarshalFrame(ctx context.Context, f *api.Frame) ([]byte, xerror.IError) {
	frameHeadLength := binary.Size(f.FrameHead)

	// 加密header
	temp, err := json.Marshal(&f.Header)
	if err != nil {
		s.ILogger.Error(ctx, "json marshal Header %+v fail %v", f.Header, err)
		return nil, xerror.Extend(xerror.ErrInternalError, err.Error())
	}
	encryptHeader, xerr := s.Encode(ctx, int(f.Encryption), temp)
	if err != nil {
		s.ILogger.Error(ctx, "encrypt header %s fail %v", temp, err)
		return nil, xerr
	}
	// 加密data
	encryptData, xerr := s.Encode(ctx, int(f.Encryption), f.Data)
	if err != nil {
		s.ILogger.Error(ctx, "encrypt data fail %v", err)
		return nil, xerr
	}
	headerLength := len(encryptHeader)
	dataLength := len(encryptData)

	contentLength := frameHeadLength
	contentLength += binary.Size(TLVHeader{})
	contentLength += headerLength
	contentLength += binary.Size(TLVHeader{})
	contentLength += dataLength

	// buffer
	contentBuffer := bytes.NewBuffer(make([]byte, 0, contentLength))

	// 帧信息
	err = binary.Write(contentBuffer, binary.BigEndian, f.FrameHead)
	if err != nil {
		s.ILogger.Error(ctx, "binary write frame head %+v fail %v\n", f.FrameHead, err)
		return nil, xerror.Extend(xerror.ErrInternalError, err.Error())
	}

	// header
	err = binary.Write(contentBuffer, binary.BigEndian, uint16(TLVTypeHeader))
	if err != nil {
		s.ILogger.Error(ctx, "binary write Header type %+v fail %v\n", TLVTypeHeader, err)
		return nil, xerror.Extend(xerror.ErrInternalError, err.Error())
	}
	err = binary.Write(contentBuffer, binary.BigEndian, uint32(headerLength))
	if err != nil {
		s.ILogger.Error(ctx, "binary write Header length %+v fail %v\n", headerLength, err)
		return nil, xerror.Extend(xerror.ErrInternalError, err.Error())
	}
	err = binary.Write(contentBuffer, binary.BigEndian, encryptHeader)
	if err != nil {
		s.ILogger.Error(ctx, "binary write Header content %s fail %v\n", temp, err)
		return nil, xerror.Extend(xerror.ErrInternalError, err.Error())
	}

	// 内容
	err = binary.Write(contentBuffer, binary.BigEndian, uint16(TLVTypeUserData))
	if err != nil {
		s.ILogger.Error(ctx, "binary write data type %+v fail %v\n", TLVTypeUserData, err)
		return nil, xerror.Extend(xerror.ErrInternalError, err.Error())
	}
	err = binary.Write(contentBuffer, binary.BigEndian, uint32(dataLength))
	if err != nil {
		s.ILogger.Error(ctx, "binary write data length %+v fail %v\n", dataLength, err)
		return nil, xerror.Extend(xerror.ErrInternalError, err.Error())
	}
	err = binary.Write(contentBuffer, binary.BigEndian, encryptData)
	if err != nil {
		s.ILogger.Error(ctx, "binary write data fail %v\n", err)
		return nil, xerror.Extend(xerror.ErrInternalError, err.Error())
	}

	// crc
	content := contentBuffer.Bytes()
	checkSum, xerr := calculateCheckSum(content[frameHeadLength:])
	if err != nil {
		s.ILogger.Error(ctx, "calculate crc check sum fail %v", err)
		return nil, xerr
	}
	copy(content[frameHeadLength-4:frameHeadLength], checkSum)

	return content, nil
}

func (s *service) UnmarshalFrame(ctx context.Context, frame []byte, f *api.Frame) xerror.IError {
	frameHeadLength := binary.Size(f.FrameHead)
	headerLength := binary.Size(f.Header)

	minLength := frameHeadLength + 2 + 4 + headerLength + 2 + 4
	if len(frame) < minLength {
		s.ILogger.Error(ctx, "frame length %+v less than %+v\n", len(frame), minLength)
		return xerror.Extend(xerror.ErrInvalidParam, "frame length less than min length")
	}

	contentReader := bytes.NewReader(frame)

	err := binary.Read(contentReader, binary.BigEndian, &f.FrameHead)
	if err != nil {
		s.ILogger.Error(ctx, "binary read frame head fail %v\n", err)
		return xerror.Extend(xerror.ErrInternalError, err.Error())
	}

	sum := CheckSum(frame[frameHeadLength:])
	if f.FrameHead.CheckSum != sum {
		s.ILogger.Error(ctx, "crc check sum fail %d != %d\n", f.FrameHead.CheckSum, sum)
		return xerror.Extend(xerror.ErrInvalidParam, "crc check sum fail")
	}

	for {
		if contentReader.Len() == 0 {
			break
		}

		tlvHeader := TLVHeader{}
		err = binary.Read(contentReader, binary.BigEndian, &tlvHeader)
		if err != nil {
			s.ILogger.Error(ctx, "binary read tlv Header fail %v\n", err)
			return xerror.Extend(xerror.ErrInternalError, err.Error())
		}

		position, err := contentReader.Seek(0, io.SeekCurrent)
		if err != nil {
			s.ILogger.Error(ctx, "seek current position fail %v", err)
			return xerror.Extend(xerror.ErrInternalError, err.Error())
		}
		if position+int64(tlvHeader.Length) > int64(len(frame)) {
			s.ILogger.Error(ctx, "frame tlv not complete\n")
			return xerror.Extend(xerror.ErrInvalidParam, "frame tlv not complete")
		}

		content := frame[position : position+int64(tlvHeader.Length)]
		content, err = s.Decode(ctx, int(f.Encryption), content)
		if err != nil {
			s.ILogger.Error(ctx, "decrypt tlv content fail %v", err)
			return xerror.Extend(xerror.ErrInternalError, err.Error())
		}

		switch tlvHeader.Type {
		case TLVTypeHeader:
			err = json.Unmarshal(content, &f.Header)
			if err != nil {
				s.ILogger.Error(ctx, "json unmarshal Header %s fail %v", content, err)
				return xerror.Extend(xerror.ErrInternalError, err.Error())
			}
		case TLVTypeUserData:
			f.Data = content
		}

		_, err = contentReader.Seek(int64(tlvHeader.Length), io.SeekCurrent)
		if err != nil {
			s.ILogger.Error(ctx, "seek to next tlv fail %v", err)
			return xerror.Extend(xerror.ErrInternalError, err.Error())
		}
	}
	return nil
}
