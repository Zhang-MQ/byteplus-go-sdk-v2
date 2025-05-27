# 集成SDK

在调用接口时，推荐在项目中集成 SDK 的方式进行接入。通过使用 SDK，不仅可以简化开发流程、加快功能集成速度，还能有效降低后期的维护成本。火山引擎 SDK 的集成主要包括以下三个步骤：引入 SDK、配置访问凭证，以及编写接口调用代码。本文将结合常见使用场景，详细说明各步骤的实现方法及注意事项。

# 环境要求

1. Go版本最低1.14+，如使用方舟服务(service/arkruntime)，需要使用go1.18+
2. 建议是用go mod的方式进行包管理


# 安全设置访问凭据

为了防止访问凭据泄露，建议不要在代码中以明文形式硬编码访问凭据。火山引擎提供了多种安全的凭据加载方式，比如将凭据存储到环境变量中。

## 环境变量设置

### Linux 设置

> ⚠️ **注意事项**
>
> **使用export命令配置的临时环境变量仅当前会话有效，当会话退出之后所设置的环境变量将会丢失。若需长期保留环境变量，可将export命令配置到对应操作系统的启动配置文件中**。


| Key                                                   | 命令                                                    |
| ----------------------------------------------------- |-------------------------------------------------------|
| BYTEPLUS_ACCESS_KEY_ID（或BYTEPLUS_ACCESS_KEY）     | export BYTEPLUS_ACCESS_KEY_ID=yourAccessKeyID         |
| BYTEPLUS_SECRET_ACCESS_KEY（或BYTEPLUS_SECRET_KEY） | export BYTEPLUS_SECRET_ACCESS_KEY=yourSecretAccessKey |
| BYTEPLUS_SESSION_TOKEN                               | export BYTEPLUS_SESSION_TOKEN=yourSessionToken       |

**验证是否设置成功：** 执行echo $BYTEPLUS_ACCESS_KEY_ID命令，如果返回正确的AccessKey ID，则说明配置成功。

### Windows 设置

下面提供了**图形化界面设置**和**命令行设置**两种方式。
**验证是否设置成功：**
单击**开始**（或快捷键：**Win+R**）> **运行**（输入 cmd）> **确定**（或按 Enter 键），打开命令提示符，执行echo %BYTEPLUS_ACCESS_KEY_ID%、echo %BYTEPLUS_SECRET_ACCESS_KEY%、echo %BYTEPLUS_SESSION_TOKEN%命令。若返回正确的值，则说明配置成功。

#### 图形化界面设置

以下为Windows 10中通过图形用户界面设置环境变量的步骤。
在桌面右键单击**此电脑**，选择**属性>高级系统设置>环境变量>系统变量/用户变量>新建**，完成以下配置：


| 变量             | 示例                                                   |
| ---------------- | ------------------------------------------------------ |
| AccessKey Id     | 变量名：BYTEPLUS_ACCESS_KEY_ID<br />变量值：*****     |
| AccessKey Secret | 变量名：BYTEPLUS_SECRET_ACCESS_KEY<br />变量值：***** |
| Session Token    | 变量名：BYTEPLUS_SESSION_TOKEN<br />变量值：*****     |

#### 命令行设置

以管理员身份打开命令提示符，并使用以下命令在系统中新增环境变量。

复制  
setx BYTEPLUS_ACCESS_KEY_ID yourAccessKeyID /M  
setx BYTEPLUS_SECRET_ACCESS_KEY yourAccessKeySecret /M  
setx BYTEPLUS_SESSION_TOKEN yourSessionToken /M

> ⚠️ 注意事项
>
> 其中/M表示系统级环境变量，设置用户级环境变量时可以不携带该参数。

# 访问凭据

为保障资源访问安全，火山引擎 SDK 支持三种主流的认证方式：**AK/SK**、**STS 临时凭证** 和 **AssumeRole**。不同认证方式适用于不同场景，开发者可根据业务需求选择合适的方式接入。

环境变量设置可以参考这里:[**环境变量设置**](#环境变量设置)

## AK、SK设置

AK/SK 是由火山引擎用户在控制台创建的一对永久访问密钥。SDK 使用该密钥对每次请求进行签名，从而完成身份验证。

> ⚠️ 注意事项
>
> 1. 不得在客户端嵌入或暴露 AK/SK。
> 2. 推荐使用配置中心或环境变量存储密钥。
> 3. 配置合理的最小权限访问策略。

```go
func main() {
    ak, sk, region := "Your AK", "Your SK", "cn-beijing"
    config := byteplus.NewConfig().
       WithRegion(region).
       // 1. credentials.NewStaticCredentials 是输入静态ak和sk可能泄漏会导致AK/SK泄漏，生产环境不能这样使用
       WithCredentials(credentials.NewStaticCredentials(ak, sk, ""))
       // 2. credentials.NewEnvCredentials() 不用传入任何参数，会从环境变量中读取：BYTEPLUS_ACCESS_KEY_ID、BYTEPLUS_SECRET_ACCESS_KEY、BYTEPLUS_SESSION_TOKEN，生产环境建议使用这个
       // WithCredentials(credentials.NewEnvCredentials())

    sess, err := session.NewSession(config)
    if err != nil {
       panic(err)
    }
}
```

## STS Token设置

STS（Security Token Service）是火山引擎提供的临时访问凭证机制。开发者通过服务端调用 STS 接口获取临时凭证（临时 AK、SK 和 Token），有效期可配置，适用于安全要求较高的场景。

> ⚠️ 注意事项
>
> 1. 最小权限： 仅授予调用方访问所需资源的最小权限，避免使用 * 通配符授予全资源、全操作权限。
> 2. 设置合理的有效期: 请根据实际情况设置合理有效期，越短越安全，建议不要超过1小时。

```Go
func main() {
    ak, sk,token,region := "Your AK", "Your SK", "Your token", "cn-beijing"
    config := byteplus.NewConfig().
       WithRegion(region).
       // 1. credentials.NewStaticCredentials 是输入静态ak和sk可能泄漏会导致AK/SK泄漏，生产环境不能这样使用
       WithCredentials(credentials.NewStaticCredentials(ak, sk, token))
       // 2. credentials.NewEnvCredentials() 不用传入任何参数，会从环境变量中读取：BYTEPLUS_ACCESS_KEY_ID、BYTEPLUS_SECRET_ACCESS_KEY、BYTEPLUS_SESSION_TOKEN，生产环境建议使用这个
       // WithCredentials(credentials.NewEnvCredentials())
    sess, err := session.NewSession(config)
    if err != nil {
       panic(err)
    }
}
```

## AssumeRole

动态访问凭证信息，支持动态刷新，在STS临时Token过期前60S会进行自动的刷新，避免临界时间点Token过期

> ⚠️ 注意事项
>
> 1. 最小权限： 仅授予调用方访问所需资源的最小权限，避免使用 * 通配符授予全资源、全操作权限。
> 2. 设置合理的有效期: 请根据实际情况设置合理有效期，越短越安全，最长不能超过12小时。
> 3. 细粒度角色: 角色应绑定精细的访问控制策略，仅允许访问特定服务、资源、操作，防止角色滥用。

```go
func main() {
    ak, sk, region := "Your AK", "Your SK", "cn-beijing"
    config := byteplus.NewConfig().
        WithRegion(region).
        WithCredentials(credentials.NewStsCredentials(credentials.StsValue{
            AccessKey:  ak,         // 子账号AK,最好从环境变量获取：os.Getenv("BYTEPLUS_ACCESS_KEY_ID")
            SecurityKey: sk,        // 子账号SK，最好从环境变量获取：os.Getenv("BYTEPLUS_SECRET_ACCESS_KEY")
            RoleName:   "RoleName", // 扮演角色名称
            Host:       "Host",     // 请求的sts域名
            Region:     "Region",   // 请求sts的region信息
            AccountId:  "123456",   // 被扮演的主账号ID，即角色所属的主账号ID
            Schema:     "Schema",   // 请求sts的schema信息
            Timeout:    5 * time.Second, // 请求sts的超时时间
            DurationSeconds: 900,        // STS临时凭证过期时长，单位为秒
        }))
 

    sess, err := session.NewSession(config)
    if err != nil {
        panic(err)
    }
}
```

# EndPoint配置

## 自定义Endpoint

> - **默认**  
>   open.volcengineapi.com

用户可以通过在初始化客户端时指定Endpoint

```go
func main() {
    region := "cn-beijing"
    config := byteplus.NewConfig().
       WithCredentials(credentials.NewEnvCredentials()). //环境变量配置：BYTEPLUS_ACCESS_KEY_ID、BYTEPLUS_SECRET_ACCESS_KEY、BYTEPLUS_SESSION_TOKEN
       WithRegion(region).
       // 自定义Endpoint
       WithEndpoint("ecs..volcengineapi.com")
     sess, err := session.NewSession(config)
     if err != nil {
        panic(err)
     }
}
```

## 自定义RegionId

```go
func main() {
    regionId := "cn-beijing"
    config := byteplus.NewConfig().
       WithCredentials(credentials.NewEnvCredentials()). //环境变量配置：BYTEPLUS_ACCESS_KEY_ID、BYTEPLUS_SECRET_ACCESS_KEY、BYTEPLUS_SESSION_TOKEN
       WithRegion(regionId). // 自定义regionId
       WithEndpoint("ecs.cn-beijing-autodriving.volcengineapi.com")
	sess, err := session.NewSession(config)
        if err != nil {
        panic(err)
        }
}
```

## 自动化Endpoint寻址

> - **默认**  
> 默认支持自动寻址，无需手动指定Endpoint

为了简化用户配置，Vocoengine 提供了灵活的 Endpoint 自动寻址机制。用户无需手动指定服务地址，SDK 会根据服务名称、区域（Region）等信息自动拼接出合理的访问地址，并支持用户自定义DualStack（双栈）支持。
### Endpoint默认寻址
**Endpoint默认寻址逻辑**
1. 是否自举Region  
   SDK 仅对部分预设区域（如 ap-southeast-2、ap-southeast-3）或用户配置的区域执行自动寻址；其他区域默认返回Endpoint：open.byteplusapi.com。  
   用户可通过环境变量 BYTEPLUS_BOOTSTRAP_REGION_LIST_CONF 或代码中自定义 customBootstrapRegion 来扩展控制区域列表。
2. DualStack 支持（IPv6）  
   SDK 支持双栈网络（IPv4 + IPv6）访问地址，自动启用条件如下：    
   显式传入参数 useDualStack = true，或设置环境变量 BYTEPLUS_ENABLE_DUALSTACK=true，优先级useDualStack>BYTEPLUS_ENABLE_DUALSTACK  
   启用后，域名后缀将从 byteplusapi.com 切换为 byteplus-api.com。
3. 根据服务名和区域自动构造 Endpoint 地址，规则如下：  
   **全局服务（如 BILLING、IAM）**  
   使用 <服务名>.byteplusapi.com（或启用双栈时使用 byteplus-api.com）。  
   示例：billing.byteplusapi.com  
   **区域服务（如 ECS、VPC）**  
   使用 <服务名>.<区域名>.byteplusapi.com 作为默认 Endpoint。  
   示例：ecs.cn-beijing.byteplusapi.com

**代码示例：**
```go
func main() {
    regionId := "cn-beijing"
    config := byteplus.NewConfig().
        WithCredentials(credentials.NewEnvCredentials()). //环境变量配置：VOLCSTACK_ACCESS_KEY_ID、VOLCSTACK_SECRET_ACCESS_KEY、VOLCSTACK_SESSION_TOKEN
        WithRegion(regionId).
        WithUseDualStack(true). // 定义是否启用双栈网络（IPv4 + IPv6）访问地址，默认false；也可以使用环境变量VOLC_ENABLE_DUALSTACK=true
        WithBootstrapRegion(map[string]struct{}{
        "cn-beijing-autodriving": {},
        "cn-shanghai-autodriving": {},
        }) // 自定义自举Region；也可以使用环境变量VOLC_BOOTSTRAP_REGION_LIST_CONF	
    sess, err := session.NewSession(config)
    if err != nil {
        panic(err)
    }
}
```

# Http连接池配置

> - **默认**  
>   最大空闲连接数（MaxIdleConns） - 100  
>   空闲连接存活时间（IdleConnTimeout） - 90  
>   每个路由最大连接数（MaxIdleConnsPerHost） - 2  

最大空闲连接数、空闲连接存活时间、每个路由最大连接数没有直接提供参数设置，可以通过自定义HTTPClient实现

```go
func main() {
    region := "cn-beijing"
    transport := &http.Transport{
       Proxy: http.ProxyFromEnvironment,
       DialContext: (&net.Dialer{
          Timeout:   30 * time.Second, 
          KeepAlive: 30 * time.Second,
          DualStack: true,
       }).DialContext,
       MaxIdleConns:          100,   // 所有host加起来的最大空闲连接数
       IdleConnTimeout:       90 * time.Second,  // 空闲连接最大存活时间
       MaxIdleConnsPerHost:   10,     // 每个host（路由）最大空闲连接数
       TLSHandshakeTimeout:   10 * time.Second,
       ExpectContinueTimeout: 1 * time.Second,
    }

    client := &http.Client{
       Transport: transport,
       Timeout:   60 * time.Second, // 这个相当于ReadTimeout
    }
    config := byteplus.NewConfig().
       WithRegion(region).
       WithHTTPClient(client). 
       WithCredentials(credentials.NewEnvCredentials()) //环境变量配置：BYTEPLUS_ACCESS_KEY_ID、BYTEPLUS_SECRET_ACCESS_KEY、BYTEPLUS_SESSION_TOKEN

    sess, err := session.NewSession(config)
    if err != nil {
       panic(err)
    }
    svc := ecs.New(sess)
}
```

# Https请求配置

## 指定Shema

> - **默认**  
>   https

schema是参数DisableSSL来控的，为true表示schema为http，为false表示schema为https；建议使用HTTPS，这样可以提升数据传输的安全性。若不设置，则使用默认支持的请求协议类型（HTTPS）

```go
func main() {
    region := "cn-beijing"
    config := byteplus.NewConfig().
       WithRegion(region).
       WithDisableSSL(true). //true 表示schema为http，false表示为https，默认为false
       WithCredentials(credentials.NewEnvCredentials()) // 环境变量配置：BYTEPLUS_ACCESS_KEY_ID、BYTEPLUS_SECRET_ACCESS_KEY、BYTEPLUS_SESSION_TOKEN

    sess, err := session.NewSession(config)
    if err != nil {
       panic(err)
    }
    svc := ecs.New(sess)
}
```

## 忽略SSL验证

> - **默认**  
>   不忽略（开启ssl认证）

没有直接提供参数设置，可以通过自定义HTTPClient实现

```go
func main() {
    region :=  "cn-beijing"
    transport := &http.Transport{
       Proxy: http.ProxyFromEnvironment,
       DialContext: (&net.Dialer{
          Timeout:   30 * time.Second, 
          KeepAlive: 30 * time.Second,
          DualStack: true,
       }).DialContext,
       MaxIdleConns:          100,   // 所有host加起来的最大空闲连接数
       IdleConnTimeout:       90 * time.Second,  // 空闲连接最大存活时间
       MaxIdleConnsPerHost:   10,     // 每个host（路由）最大空闲连接数
       TLSHandshakeTimeout:   10 * time.Second,
       ExpectContinueTimeout: 1 * time.Second,
       TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, //跳过服务器证书校验，即使证书无效也继续访问
    }

    client := &http.Client{
       Transport: transport,
       Timeout:   60 * time.Second, // 这个相当于ReadTimeout
    }
    config := byteplus.NewConfig().
       WithRegion(region).
       WithHTTPClient(client). 
       WithCredentials(credentials.NewEnvCredentials()) //环境变量配置：BYTEPLUS_ACCESS_KEY_ID、BYTEPLUS_SECRET_ACCESS_KEY、BYTEPLUS_SESSION_TOKEN

    sess, err := session.NewSession(config)
    if err != nil {
       panic(err)
    }
    svc := ecs.New(sess)
}
```

## 指定TLS协议版本

> - **默认**  
>   \>=TLS 1.2

目前只支持自定义HTTPClient的方式实现；如果没有特殊要求，建议不要修改。

```go
func main() {
    region := "cn-beijing"
    transport := &http.Transport{
       Proxy: http.ProxyFromEnvironment,
       DialContext: (&net.Dialer{
          Timeout:   30 * time.Second, 
          KeepAlive: 30 * time.Second,
          DualStack: true,
       }).DialContext,
       MaxIdleConns:          100,   // 所有host加起来的最大空闲连接数
       IdleConnTimeout:       90 * time.Second,  // 空闲连接最大存活时间
       MaxIdleConnsPerHost:   10,     // 每个host（路由）最大空闲连接数
       TLSHandshakeTimeout:   10 * time.Second,
       ExpectContinueTimeout: 1 * time.Second,
       TLSClientConfig: &tls.Config{
           MinVersion: tls.VersionTLS12, // 只允许 TLS 1.2
           MaxVersion: tls.VersionTLS13, // 最大TLS1.3
       }, 
    }

    client := &http.Client{
       Transport: transport,
       Timeout:   60 * time.Second, // 这个相当于ReadTimeout
    }
    config := byteplus.NewConfig().
       WithRegion(region).
       WithHTTPClient(client). 
       WithCredentials(credentials.NewEnvCredentials()) //环境变量配置：BYTEPLUS_ACCESS_KEY_ID、BYTEPLUS_SECRET_ACCESS_KEY、BYTEPLUS_SESSION_TOKEN

    sess, err := session.NewSession(config)
    if err != nil {
       panic(err)
    }
    svc := ecs.New(sess)
}
```

# 超时配置

## 全局超时设置（Client级别）

> - **默认**  
>   ConnectTimeOut  30s  
>   ReadTimeout       不限制  
>   备注：默认用的是http.DefaultClient  

暂不支持直接设置ConnectTimeOut和ReadTimeout配置，可以通过自定义HttpClient来实现

```Go
func main() {
    region := "cn-beijing"
    transport := &http.Transport{
       Proxy: http.ProxyFromEnvironment,
       DialContext: (&net.Dialer{
          Timeout:   30 * time.Second, // 这个相当于ConnectTimeOut
          KeepAlive: 30 * time.Second,
          DualStack: true,
       }).DialContext,
       MaxIdleConns:          100,
       IdleConnTimeout:       90 * time.Second,
       TLSHandshakeTimeout:   10 * time.Second,
       ExpectContinueTimeout: 1 * time.Second,
    }

    client := &http.Client{
       Transport: transport,
       Timeout:   60 * time.Second, // 这个相当于ReadTimeout
    }
    config := byteplus.NewConfig().
       WithRegion(region).
       WithHTTPClient(client). 
       WithCredentials(credentials.NewEnvCredentials()) // 环境变量配置：BYTEPLUS_ACCESS_KEY_ID、BYTEPLUS_SECRET_ACCESS_KEY、BYTEPLUS_SESSION_TOKEN

    sess, err := session.NewSession(config)
    if err != nil {
       panic(err)
    }
    svc := ecs.New(sess)
}
```

## 单接口指定超时设置

单个接口超时配置需要调用以**WithContext**结尾的接口

```go
func main() {
	region :=  "cn-beijing"
    config := byteplus.NewConfig().
       WithRegion(region).
       WithCredentials(credentials.NewEnvCredentials()) //环境变量配置：BYTEPLUS_ACCESS_KEY_ID、BYTEPLUS_SECRET_ACCESS_KEY、BYTEPLUS_SESSION_TOKEN
    sess, err := session.NewSession(config)
    if err != nil {
       panic(err)
    }
    svc := ecs.New(sess)
    associateInstancesIamRoleInput := &ecs.AssociateInstancesIamRoleInput{
       IamRoleName: byteplus.String("EcsTestRole"),
       InstanceIds: byteplus.StringSlice([]string{"i-3tiefmkskq3vj0******"}),
    }

    // 创建带5秒超时的上下文
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    // 防止资源泄漏
    defer cancel() 

    // 调用WithContext结尾的接口
    resp, err := svc.AssociateInstancesIamRoleWithContext(ctx, associateInstancesIamRoleInput)
    if err != nil {
       panic(err)
    }
    fmt.Println(resp)
}
```

# 重试机制

请求的处理逻辑内置了网络异常重试逻辑，即当遇到网络异常问题或限流错误时，系统会自动尝试重新发起请求，以确保服务的稳定性和可靠性。若请求因业务逻辑错误而报错，例如参数错误、资源不存在等情况，SDK将不会执行重试操作，这是因为业务层面的错误通常需要应用程序根据具体的错误信息做出相应的处理或调整，而非简单地重复尝试

## 开启重试机制

> - **默认**  
>   开启（3次）

如果想关闭，可以把最大尝试次数改为0

```go
func main() {  
        region := "cn-beijing"
        // Configure retry settings
        config := byteplus.NewConfig().
                WithRegion(region).
                WithDisableSSL(true).
                WithCredentials(credentials.NewEnvCredentials()). //环境变量配置：BYTEPLUS_ACCESS_KEY_ID、BYTEPLUS_SECRET_ACCESS_KEY、BYTEPLUS_SESSION_TOKEN
                // 关闭重试
                WithMaxRetries(0)

        sess, err := session.NewSession(config)
        if err != nil {
                panic(err)
        }
        svc := ecs.New(sess)
}
```

## 重试次数

> - **默认**  
>   3次

1. 设置默认次数

```go
func main() {
        region := "cn-beijing"
        // Configure retry settings
        config := byteplus.NewConfig().
                WithRegion(region).
                WithDisableSSL(true).
                WithCredentials(credentials.NewEnvCredentials()). //环境变量配置：BYTEPLUS_ACCESS_KEY_ID、BYTEPLUS_SECRET_ACCESS_KEY、BYTEPLUS_SESSION_TOKEN
                // 设置最大重试次数 (default is DefaultRetryerMaxNumRetries)
                WithMaxRetries(4)

        sess, err := session.NewSession(config)
        if err != nil {
                panic(err)
        }
        svc := ecs.New(sess)
}
```

## 自定义重试错误码

在调用接口的时候可以根据业务需求，自定义重试的错误码(服务端返回的错误码)。

```go
func main() {
        region := "cn-beijing"
        // Configure retry settings
        config := byteplus.NewConfig().
                WithRegion(region).
                WithDisableSSL(true).
                WithCredentials(credentials.NewEnvCredentials()). //环境变量配置：BYTEPLUS_ACCESS_KEY_ID、BYTEPLUS_SECRET_ACCESS_KEY、BYTEPLUS_SESSION_TOKEN
                // 设置最大重试次数 (default is DefaultRetryerMaxNumRetries)
                WithMaxRetries(4)

        sess, err := session.NewSession(config)
        if err != nil {
                panic(err)
        }
        svc := ecs.New(sess)
        describeAvailableResourceInput := &ecs.DescribeAvailableResourceInput{
            DestinationResource: byteplus.String("InstanceType"),
            InstanceTypeId:      byteplus.String("ecs.g2i.large"),
            ZoneId:              byteplus.String("cn-*****"),
        }
        ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
        defer cancel()
        resp, err := svc.DescribeAvailableResourceWithContext(ctx, describeAvailableResourceInput, func(request *request.Request) {
            // 自定义重试错误码
            request.RetryErrorCodes = []string{"InvalidAccessKey"}
        })
        if err != nil {
            panic(err)
        }
}
```

# 异常处理

在调用接口时，可能会返回不同类型的错误。客户可根据具体的错误类型和错误码，采取有针对性的处理策略。例如，对于网络异常可选择重试，对于业务逻辑错误则应根据错误信息进行参数调整或业务逻辑修正，从而提升系统的健壮性和用户体验。

错误分类：


| 错误类型                         | 错误描述                                                      |
| -------------------------------- | ------------------------------------------------------------- |
| 1. 配置缺失/验证错误             | 某些全局配置没有配置，如：Region,endpoint,请求body格式验证等  |
| 2. 网络错误                      | SDK 在尝试发起 HTTP 请求前就失败，或在建立连接/发送请求时失败 |
| 3. 超时错误                      | 请求已成功发出，但等待响应超时（可细分为连接超时、读取超时）  |
| 4. 业务错误                      | 请求成功到达服务器，但返回的是业务逻辑错误（参数不合法）      |
| 5. 认证错误                      | 请求签名无效、Token 过期或缺失                                |
| 6. 限流/熔断错误(ThrottingError) | 服务端返回限流标志，通常 HTTP 429，或内部系统降级             |
| 7. 反序列化数据错误              | 对json或xml反序列化报的错误                                   |

代码示例：

```go
package main

import (
    "context"
    "errors"
    "fmt"
    "github.com/volcengine/volcengine-go-sdk/service/ecs"
    "github.com/volcengine/volcengine-go-sdk/volcengine"
    "github.com/volcengine/volcengine-go-sdk/volcengine/credentials"
    "github.com/volcengine/volcengine-go-sdk/volcengine/request"
    "github.com/volcengine/volcengine-go-sdk/volcengine/session"
    "github.com/volcengine/volcengine-go-sdk/volcengine/volcengineerr"
    "net"
    "net/http"
    "time"
)

func main() {
    region := "cn-beijing"
    config := byteplus.NewConfig().
       WithRegion(region).
       WithCredentials(credentials.NewEnvCredentials())
    sess, err := session.NewSession(config)
    var be bytepluserr.Error
    if err != nil {
       // 1. 配置/参数验证错误
       if errors.As(err, &be) {
          switch be.Code() {
          case "LoadCustomCABundleError":
             fmt.Println("1. 配置缺失/验证错误：当在配置文件中未找到配置文件时会出现此错误")
          case "SharedConfigLoadError":
             fmt.Println("1. 配置缺失/验证错误：SharedConfigLoadError是共享配置文件加载失败时的错误。")
          case "SharedConfigProfileNotExistsError":
             fmt.Println("1. 配置缺失/验证错误：共享配置文件不存在错误是共享配置出现的一种错误，当在配置文件中未找到配置文件时就会出现该错误。")
          case "SharedConfigAssumeRoleError":
             fmt.Println("1. 配置缺失/验证错误：SharedConfigAssumeRoleError是共享配置中的一种错误，当配置文件包含角色假设信息，但该信息无效或不完整时会出现此错误。")
          case "CredentialRequiresARNError":
             fmt.Println("1. 配置缺失/验证错误：credential type (source_profile|credential_source|web_identity_token_file) requires role_arn, profile ")
          }

       } else {
          fmt.Println("处理其它错误")
       }
       panic(err)
    }
    svc := ecs.New(sess)

    tags := make([]*ecs.TagForCreateKeyPairInput, 0, 2)
    tags = append(tags, &ecs.TagForCreateKeyPairInput{Key: byteplus.String("testTag")})
    createKeyPairInput := &ecs.CreateKeyPairInput{
       KeyPairName: byteplus.String("testKeyPairName"),
       Tags:        tags,
    }

    _, err = svc.CreateKeyPair(createKeyPairInput)
    if err != nil {
       var requestFailure bytepluserr.RequestFailure // 服务端返回的错误
       var errInvalidParam request.ErrInvalidParam     // 参数验证错误
       var unmarshalError bytepluserr.UnmarshalError // 返回数据解析错误
       var batchedErrors bytepluserr.BatchedErrors   // 批量错误
       // 请求未达到服务前参数验证
       if errors.As(err, &errInvalidParam) {
          fmt.Println("1. 配置缺失/验证错误：", errInvalidParam.Field(), errInvalidParam.Error())
       // 请求到达服务端，服务端返回错误   
       } else if errors.As(err, &requestFailure) {
          fmt.Println("请求错误：", requestFailure.RequestID(), requestFailure.StatusCode())
          switch requestFailure.Code() {
          // 请求过于频繁
          case "FlowLimitExceeded":
             fmt.Println("6. 限流/熔断错误(ThrottingError): 请求过于频繁，超出了限速。请降低请求QPS，如果有提升限速需求")
          case "AccessDenied":
             fmt.Println("5. 认证错误：用户拥有的权限不支持当前操作")
          case "InvalidActionOrVersion":
             fmt.Println("4. 业务错误：请求接口不存在")
          case "InvalidAccessKey":
             fmt.Println("5. 认证错误：请求的Access Key不合法。请检查Access key Id和Secret Access Key是否正确，注意不要有多余的空格符号")
          case "InternalServiceTimeout":
             fmt.Println("3. 超时错误：服务内部执行超时")
          case "InvalidAuthorization":
             fmt.Println("5. 认证错误：Authorization头格式错误，构造的 Authorization Header 不正确，比如没有填 Region 字段、字符不在合法字符集中，请检查Authorization")
          case "InvalidCredential":
             fmt.Println("5. 认证错误： Authorization头中的Credential格式错误，比如 AK 不在合法字符集中检查Credential")
          case "InvalidSecretToken":
             fmt.Println("5. 认证错误：错误的STS（临时安全凭证），可能是多种错误，例如签名错误、过期等。请检查AssumeRole产生的临时凭证是否过期，以及签名是正确")
          default:
             fmt.Println("4. 业务错误：可以把其它未处理的后端返回错误归结为业务错误")
          }
        // 反序列化数据错误
       } else if errors.As(err, &unmarshalError) {
          fmt.Println("7. 反序列化数据错误：", string(unmarshalError.Bytes()))
       } else if errors.As(err, &batchedErrors) {
          for _, e := range batchedErrors.OrigErrs() {
             fmt.Println("批量错误:", e)
          }
       // 发送请求，但没有到达后端服务   
       }else if errors.As(err, &be) {
          switch be.Code() {
          case "RequestCanceled":
             fmt.Println("3. 超时错误：这里是请求接口传入context上下文超时")
          case "RequestError":
             if be.OrigErr() != nil {
                var netErr net.Error
                var dnsError *net.DNSError
                if errors.As(be.OrigErr(), &dnsError) {
                    fmt.Println("2. 网络错误：DNS解析错误处理")
                }else if errors.As(be.OrigErr(), &netErr) && netErr.Timeout() {
                   var oPError *net.OpError
                   if errors.Is(be.OrigErr(), context.DeadlineExceeded) {
                      fmt.Println("3. 超时错误：http.Client Timeout(ReadTimeout)....")
                   } else if errors.As(be.OrigErr(), &oPError) && oPError.Op == "dial" {
                      fmt.Println("3. 超时错误：http.Client Transport.Dialer Timeout(ConnectTimeout)....")
                   } else {
                      fmt.Println("3. 超时错误：其它超时处理")
                   }
                }
             }
          default:
             fmt.Println("处理其它错误")
          }
       }  else {
          fmt.Println("处理其它错误")
       }

    }

}

```

# Debug机制

为便于客户在处理请求时进行问题排查和调试，SDK 支持日志功能，并提供多种日志级别设置。客户可根据实际需求配置日志级别，获取详细的请求与响应信息，以提升排障效率和系统可 observability（可观测性）。

> 1. LogOff - 关闭调试日志(默认)
> 2. LogDebug - 开启log调试日志
>    LogDebug又细分为：
>    - LogDebugWithSigning
>      记录请求签名和预签名事件
>      用于调试请求的签名细节
>      会同时启用 LogDebug
>    - LogDebugWithHTTPBody
>      记录 HTTP 请求和响应的 body（除 headers 和路径外）
>      用于查看 SDK 请求和响应的完整内容
>      会同时启用 LogDebug
>    - LogDebugWithRequestRetries
>      记录服务请求的重试情况
>      用于跟踪服务请求何时被重试
>      会同时启用 LogDebug
>    - LogDebugWithRequestErrors
>      记录服务请求构建、发送、验证或反序列化失败的情况
>    - LogDebugWithEventStreamBody
>      记录 EventStream 请求和响应的 body
>      用于查看 EventStream 消息内容
>      会同时启用 LogDebug
>    - LogInfoWithInputAndOutput
>      记录结构体(STRUCT)的输入和输出
>      会同时启用 LogInfo
>    - LogDebugWithInputAndOutput
>      记录结构体(STRUCT)的输入和输出
>      会同时启用 LogDebug

```go
func main() {
    region :=  "cn-beijing"
    config := byteplus.NewConfig().
       WithRegion(region).
       WithLogLevel(byteplus.LogDebugWithInputAndOutput). // 设置日志打印级别，不设置默认LogOff 不打印日志
       WithCredentials(credentials.NewEnvCredentials()) // 环境变量配置：BYTEPLUS_ACCESS_KEY_ID、BYTEPLUS_SECRET_ACCESS_KEY、BYTEPLUS_SESSION_TOKEN
    sess, err := session.NewSession(config)
    if err != nil {
       panic(err)
    }
    svc := ecs.New(sess)

    tags := make([]*ecs.TagForCreateKeyPairInput, 0, 2)
    tags = append(tags, &ecs.TagForCreateKeyPairInput{Key: byteplus.String("test")})
    createKeyPairInput := &ecs.CreateKeyPairInput{
       KeyPairName: byteplus.String("test"),
       Tags:        tags,
    }

    // 复制代码运行示例，请自行打印API返回值。
    _, err = svc.CreateKeyPair(createKeyPairInput)
    if err != nil {
       // 复制代码运行示例，请自行打印API错误信息。
       panic(err)
    }
}
```

# 指定日志Logger

> - 默认  
>   volcengine/logger.go

在不指定logger的情况下默认使用的是volcengine/logger.go下的defaultLogger，核心代码如下：

```go
func NewDefaultLogger() Logger {
    return &defaultLogger{
       logger: log.New(os.Stdout, "", log.LstdFlags),
    }
}

// A defaultLogger provides a minimalistic logger satisfying the Logger interface.
type defaultLogger struct {
    logger *log.Logger
}

// Log logs the parameters to the stdlib logger. See log.Println.
func (l defaultLogger) Log(args ...interface{}) {
    l.logger.Println(args...)
}
```

## 自定义Logger

客户可根据业务需求，参考 SDK 默认日志实现，自定义 Logger 的输出方式。例如，可以自定义日志前缀、自定义日志输出目标（如控制台、文件或日志系统），以及实现日志内容脱敏处理等功能，从而更好地满足自身的运维与安全要求。

```go
// 自定义Logger实现
type myLogger struct {
        logger *log.Logger
}

func (l *myLogger) Log(args ...interface{}) {
       // 敏感信息过滤
	for i, arg := range args {
		if s, ok := arg.(string); ok {
			args[i] = strings.Replace(s, "KeyWord", "***", -1)
		}
	}
	l.logger.Println(args...)
}

func main() {
        region := "cn-beijing"
        // 定义输出文件路径
        file, _ := os.Create("ecs_test.log")
        // 同时输出到控制台和文件
	multiWriter := io.MultiWriter(os.Stdout, file) 
	customLogger := &myLogger{
                // 第1个参数定义输出目标(os.Stdout表示控制台输出，可以定义文件输出)
                // 第2个参数定义日志前缀
		logger: log.New(multiWriter, "[MyApp] ", log.LstdFlags|log.Lshortfile),
	}

        config := byteplus.NewConfig().
                WithRegion(region).
                WithLogLevel(byteplus.LogDebugWithInputAndOutput).
                WithCredentials(credentials.NewEnvCredentials()). //环境变量配置：BYTEPLUS_ACCESS_KEY_ID、BYTEPLUS_SECRET_ACCESS_KEY、BYTEPLUS_SESSION_TOKEN
                WithLogger(customLogger)  // 设置自定义Logger

        sess, err := session.NewSession(config)
        if err != nil {
                panic(err)
        }
        svc := ecs.New(sess)
}
```
