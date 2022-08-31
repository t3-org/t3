package infra

var banner = `
┌─┐┌─┐┌─┐┌─┐┌─┐
└─┐├─┘├─┤│  ├┤ 
└─┘┴  ┴ ┴└─┘└─┘%s
%s
Space Bench lab.
%s
____________________________________O/_______
                                    O\

`

// Banner returns the app's banner.
func Banner() string {
	return banner
}
