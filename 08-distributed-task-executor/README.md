# Distributed Task Executor
You are tasked with implementing a failover mechanism for executing tasks across multiple distributed nodes. Each node has a priority level, and you need to ensure tasks are executed reliably even when some nodes fail.
Your task is to implement the `ExecuteWithFailover` function that:

* Tries nodes in order of priority (lowest number first)
* If a node fails, immediately tries the next node without waiting
* If a node doesn't respond within 500ms, tries the next node but keeps the original request running
* Returns the first successful result, or all errors if all nodes fail
* Properly handles context cancellation throughout the process

Tags
`Concurrency` `Context` `Failover` 
