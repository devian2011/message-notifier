package transport

type Handler interface {
	HandleTemplate(command TemplateCommand) error
	HandleCustom(command CustomCommand) error
}
