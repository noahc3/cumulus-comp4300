package main

type CloudServer interface {
	LaunchProcess()
	StartServer()
	StopServer()
	RestartServer()
}
