###学习笔记

error处理的最佳实践：
 假设有一个项目的数据流结构为:\
  dao -> service -> controller\
  假设此时在dao层报error，那么此时的error在各层的传递和处理方式如下：\
  1. dao层使用 errors.Wrap(err, "") 或者 errors.WithStack(err) 将 该err携带报错时的堆栈信息返回给上层
  2. service层收到dao层的err后，有两种处理方式：\
    a. 直接返回给其上层或者使用 errors.WithMessage(err, "") 将该err携带上下文信息返回给其上层。\
    b. 在本层处理该err（处理方式包括打日志或者其他业务逻辑来处理），无需将其返回给上层。
  3. controller层收到service层报错的err后，需要将该err通过日志的方式打印出来(log.Info(fmt.Sprintf("err info with stack is: %+v", err))), 然后再做其他业务逻辑处理。
  
  
 整个err返回流可以概括为以下几个注意点：
 1. 在error第一次出现时，需要将该error及其堆栈信息（使用errors.Wrap(err, "") 或者 errors.WithStack(err)）返回给上层。
 2. 然后其上层接收到该error后，可直接继续向上返回该error或者使用 errors.WithMessage(err, "") 方法将其反悔至其上层。
 3. 顶层过在接收到err时，需要通过(log.Info(fmt.Sprintf("err info with stack is: %+v", err)))的方式进行日志打印（注意同一error全局只有这一次日志打印error）。
 4. 如果时共用的公用函数，其在返回err时，不推荐携带stack信息，直接 return error 即可。
 
  
  
  