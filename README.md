# panic
1. 在程序启动的时候，如果有强依赖的服务出现故障时`panic`退出（✅）
2. 在程序启动的时候，如果发现有配置明显不符合要求,可以`panic`退出（✅）
3. 其他情况下只要不是不可恢复的程序错误，都不应该直接`panic`应该返回`error`（✅）
4. 在程序入口处，例如`gin`中间件需要使用`recover`预防`panic`程序退出（💭虽然HTTP服务端对所有的处理函数都会使用`recover`处理`panic`，但David Symonds在GitHub的评论中 https://oreil.ly/BGOmg 提到这其实是一个决策错误）
5. 在程序中我们应该避免使用野生的`goroutine`（💡，含有io的场景是考虑协程的好时机）
    - 如果是在请求中需要执行异步任务，应该使用异步`worker`，消息通知的方式进行处理，避免请求量大时大量`goroutine`创建（💭，Go程序本来可以同时生成数百、数千甚至数万个`goroutine`）
    - 如果需要使用`goroutine`时，应该使用同一的`Go`函数进行创建，这个函数中会进行`recover`，避免因为野生`goroutine` `panic`导致主进程退出（✅，不过这个写法只支持`func()`这个类型的参数）
    
    ```go
    func Go(f func()){
        go func(){
            defer func(){
                if err := recover(); err != nil {
                    log.Printf("panic: %+v", err)
                }
            }()

            f()
        }()
    }
    ```

# error
1. 我们在应用程序中使用`github.com/pkg/errors`处理应用错误，**注意在公共库当中，我们一般不使用这个**（❓，为什么？）
2. `error`应该是函数的最后一个返回值，当`error`不为`nil`时，函数的其他返回值是不可用的状态，不应该对其他返回值做任何期待（✅，`func f()(io.Reader, *S1, error)` 在这里，我们不知道`io.Reader`中是否有数据，可能有，也有可能有一部分）
3. 错误处理的时候应该先判断错误，`if err != nil`出现错误及时返回，使代码是一条流畅的直线，避免过多的嵌套.(✅)
4. 在应用程序中出现错误时，使用`errors.New`或者`errors.Errorf`返回错误(✅)
    ```go
    func (u *usecese) usecase1() error {
        money := u.repo.getMoney(uid)
        if money < 10 {
            errors.Errorf("用户余额不足, uid: %d, money: %d", uid, money)
        }
        // 其他逻辑
        return nil
    }    
    ```    
5. 如果是调用应用程序的其他函数出现错误，请直接返回，如果需要携带信息，请使用`errors.WithMessage`(✅)
    ```go
    func (u *usecese) usecase2() error {
        name, err := u.repo.getUserName(uid)
        if err != nil {
            return errors.WithMessage(err, "其他附加信息")
        }

        // 其他逻辑
        return nil
    }  
    ```    
6. 如果是调用其他库（标准库、企业公共库、开源第三方库等）获取到错误时，请使用`errors.Wrap`添加堆栈信息(✅)
    - 切记，不要每个地方都是用`errors.Wrap`只需要在错误第一次出现时进行`errors.Wrap`即可(✅)
    - 根据场景进行判断是否需要将其他库的原始错误吞掉，例如可以把`repository`层的数据库相关错误吞掉，返回业务错误码，避免后续我们分割微服务或者更换`ORM`库时需要去修改上层代码(✅)
    - 注意我们在基础库，被大量引入的第三方库编写时一般不使用`errors.Wrap`避免堆栈信息重复(❓，没太get到)
        ```go
        func (u *usecese) usecase2() error {
            name, err := u.repo.getUserName(uid)
            if err != nil {
                return errors.WithMessage(err, "其他附加信息")
            }

            // 其他逻辑
            return nil
        }   
        ```
7. **禁止**每个出错的地方都打日志，只需要在进程的最开始的地方使用`%+v`进行统一打印，例如`http/rpc`服务的中间件
8. 错误判断使用`errors.Is`进行比较(✅)
    ```go
    func f() error {
        err := A()
        if errors.Is(err, io.EOF){
            return nil
        }

        // 其他逻辑
        return nil
    } 
    ```
9. 错误类型判断，使用`errors.As`进行赋值
    ```go
    func f() error {
        err := A()

        var errA errorA
        if errors.As(err, &errA){
            // ...
        }

        // 其他逻辑
        return nil
    }
    ```
