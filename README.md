# dgen

## 简介

drpc 是一个类似于grpc的跨语言RPC框架，dgen是其编译器部分。

本项目设计了一种简单的中立接口描述语言(IDL)，用户可以通过dgen从IDL生成目标语言的桩代码，并与drpc框架的接口进行对接。


## IDL 语法

IDL 语法类似于Protobuf，但更为简单。只包含四个关键字:

+ `enum`：用于定义枚举类型
+ `message`：用于定义复合类型
+ `service`: 用于定义服务集合
+ `optional`: 用于定义message的成员为可选（即该成员值可以为空），注：message中每个成员默认是必选的。

**基础类型**：`uint8`、`uint16`、`uint32`、`uint64`、`int8`、`int16`、`int32`、`int64`、`string`

**示例**
```protobuf
# comment
enum fruit {
    apple,
    banana
}

message HelloRequest {
    seq=1 string name;
}

message HelloResponse {
    seq=1 string name;
    optional seq=2 string reply;
}

service Greeter {
    SayHello(HelloRequest) return (HelloResponse);
    OrderFruit(fruit);
}
```

## 使用方法
```
Usage of hgen:
    -e string 
        the serialization method of message (default "", represent adopt the project's default serialization method, optional "json")
    -f string
        the path of IDL file
    -o string
    	the dirpath where the generated source code files will be placed (default ".")
    -l string
    	the target languege the IDL will be compliled
```