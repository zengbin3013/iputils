package iputils

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

// 生成指定数量的重复字符串
func generateRepeatString(origin string,count int)string{
	strCount:=strconv.Itoa(count)
	format:="%"+strCount+"s"
	dstString:=	fmt.Sprintf(format,"")
	return strings.Replace(dstString," ",origin,-1)
}

func ip2binaryIP(ip string)(string,error){
	ipSegments:=strings.Split(ip,".")
	if len(ipSegments)!=4{
		return "",errors.New("invalid ipv4 format")
	}
	ipBinaryString:=""
	for _,seg:=range ipSegments{
		intSeq,err:=strconv.Atoi(seg)
		if err!=nil{
			return "",errors.New("invalid ipv4 format")
		}
		if intSeq>255 ||intSeq<0{
			return "",errors.New("invalid ipv4 format")
		}
		ipBinaryString+=fmt.Sprintf("%08s",fmt.Sprintf("%b",intSeq))
	}
	return ipBinaryString,nil
}

func binaryIP2IPV4(binaryIP string)(string,error){
	if len(binaryIP)!=32{
		return "",errors.New("invalid ip binary format")
	}
	segment1:=binaryIP[0:8]
	segment2:=binaryIP[8:16]
	segment3:=binaryIP[16:24]
	segment4:=binaryIP[24:]
	intSegment1,err:=convBinary2Int64(segment1)
	if err!=nil{
		return "",errors.New("invalid ip binary format")
	}
	intSegment2,err:=convBinary2Int64(segment2)
	if err!=nil{
		return "",errors.New("invalid ip binary format")
	}
	intSegment3,err:=convBinary2Int64(segment3)
	if err!=nil{
		return "",errors.New("invalid ip binary format")
	}
	intSegment4,err:=convBinary2Int64(segment4)
	if err!=nil{
		return "",errors.New("invalid ip binary format")
	}
	if intSegment1>255 ||intSegment1<0 || intSegment2>255 ||intSegment2<0 ||intSegment3>255 ||intSegment3<0||intSegment4>255 ||intSegment4<0{
		return  "",errors.New("invalid ip binary format")
	}
	return fmt.Sprintf("%d.%d.%d.%d",intSegment1,intSegment2,intSegment3,intSegment4),nil
}

func convBinary2Int64(binary string)(int64,error){
	ret:=int64(0)
	binary=strings.TrimLeft(binary,"0")
	for index,c:=range binary{
		if string(c)!="0"&&string(c)!="1"{
			return 0,errors.New("invalid binary string")
		}
		intC,_:=strconv.Atoi(string(c))
		ret+=int64(intC)*int64(math.Pow(float64(2),float64(len(binary)-index-1)))
	}
	return ret,nil
}

func IPv4IsValid(ip string)bool{
	ipSegments:=strings.Split(ip,".")
	if len(ipSegments)!=4{
		return false
	}
	for _,segment:=range ipSegments{
		intSegment,err:=strconv.Atoi(segment)
		if err!=nil{
			return false
		}
		if intSegment>255 || intSegment<0{
			return false
		}
	}
	return true
}

func CIDRIsValid(cidr string)bool{
	cidrSegments:=strings.Split(cidr,"/")
	if len(cidrSegments)!=2{
		return false
	}
	if !IPv4IsValid(cidrSegments[0]){
		return false
	}
	mask,err:=strconv.Atoi(cidrSegments[1])
	if err!=nil{
		return false
	}
	if mask<0 || mask>32{
		return false
	}
	return true
}

func IPV42Int64(ip string)(int64,error){
	binaryString,err:=ip2binaryIP(ip)
	if err!=nil{
		return 0,err
	}
	ret,err:=convBinary2Int64(binaryString)
	if err!=nil{
		return 0,err
	}
	return ret,nil
}

func CIDR2IPRange(cidr string)(startIP string,endIP string ,err error){
	subnetSegments:=strings.Split(cidr,"/")
	if len(subnetSegments)!=2{
		return "","",errors.New("invalid ip subnet format,should be: 0.0.0.0/0")
	}
	binaryIP,err:=ip2binaryIP(subnetSegments[0])
	if err!=nil{
		return "","",err
	}
	mask,err:=strconv.Atoi(subnetSegments[1])
	if err!=nil{
		return "","",errors.New("invalid mask number")
	}
	if mask>32||mask<0{
		return "","",errors.New("invalid mask number")
	}
	startBinaryIP:=binaryIP[0:mask]+generateRepeatString("0",32-mask)
	endBinaryIP:=binaryIP[0:mask]+generateRepeatString("1",32-mask)
	startIP,err=binaryIP2IPV4(startBinaryIP)
	if err!=nil{
		return "","",errors.New("invalid ip format")
	}
	endIP,err=binaryIP2IPV4(endBinaryIP)
	if err!=nil{
		return "","",errors.New("invalid ip format")
	}
	return startIP,endIP,nil
}

func CheckIPInCIDR(ip string,cidr string)(bool){
	int64IP,err:=IPV42Int64(ip)
	if err!=nil{
		return false
	}

	startIP,endIP,err:=CIDR2IPRange(cidr)
	if err!=nil{
		return false
	}

	int64StartIP,err:=IPV42Int64(startIP)
	if err!=nil{
		return false
	}

	int64EndIP,err:=IPV42Int64(endIP)
	if err!=nil{
		return false
	}
	if int64IP>=int64StartIP && int64IP<=int64EndIP{
		return true
	}
	return false
}