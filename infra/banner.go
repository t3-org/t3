package infra

var banner = `
┌─┐┌─┐┌─┐┌─┐┌─┐
└─┐├─┘├─┤│  ├┤ 
└─┘┴  ┴ ┴└─┘└─┘%s
%s
T3 server
%s
____________________________________O/_______
                                    O\

`

// Banner returns the app's banner.
func Banner() string {
	return banner
}
