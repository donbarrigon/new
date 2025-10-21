package seed

type Seeds []func()

func Run() Seeds {
	return []func(){
		UserSeeder,
	}
}
