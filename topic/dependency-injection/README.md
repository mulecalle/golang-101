# golang-di-demo

https://medium.com/@mrcyna/introduction-to-dependency-injection-in-golang-e4152c874b37

“Code to interfaces, not implementations”

You need to write your program (function, class, package, service, test, module, config …) in a way that is able to accept a **behavior** (Abstraction, Interfaces) instead of Implementation (Concretion) when it comes to being dependent on other things.

In a much simpler sentence: You should be dependent on a behavior rather than a specific object/type. In our code Newsletter service is dependent on Mailgun, whereas it should be dependent on any object/type that provides mailing behavior!

In Golang the interface is a custom type that is used to specify a set of one or more method signatures and the interface is abstract. You are not defining that explicitly and Golang will detect that for you.