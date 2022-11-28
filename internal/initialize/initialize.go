package initialize

type Initializer interface {
	CheckAdminAccess() error             // Check if the user has admin access
	CheckRequirement() ([]string, error) // CheckRequirement checks the requirement for the system, returns a list of missing dependencies
	InstallDeps([]string) error          // InstallDeps installs the dependencies
	CheckDB() error                      // CheckDB checks the database access
	ConfigureDB() error                  // ConfigureDB configures the database
	checkPyenv() error                   // checkPyenv checks if pyenv is installed
	InstallPyenv() error                 // InstallPyenv installs pyenv
	ConfigurePyenv() error               // ConfigurePyenv configures pyenv
}

type DefaultInitializer struct {
	Shell string
}
