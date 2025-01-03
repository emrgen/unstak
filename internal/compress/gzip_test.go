package compress

import (
	"github.com/sirupsen/logrus"
	"testing"
)

func Test_GZipEncode(t *testing.T) {
	br := NewGZip()
	data := []byte(`Lorem ipsum dolor sit amet, consectetur adipiscing elit. Morbi vulputate mi ac lectus laoreet vehicula. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia curae; Fusce ultricies vestibulum dapibus. Nunc ac leo nibh. Nullam vitae lacinia arcu. Donec at lectus turpis. Integer finibus, justo vel maximus eleifend, tortor eros accumsan tortor, ut varius purus diam sit amet quam. Fusce finibus ac sem nec malesuada. Sed ex sapien, auctor id turpis ac, faucibus luctus felis. Maecenas vel nulla tortor. Duis nec pretium ante, vel laoreet quam. Duis facilisis consectetur turpis non pharetra. In sit amet leo et erat cursus facilisis. Sed sed malesuada leo, a dapibus justo.
Class aptent taciti sociosqu ad litora torquent per conubia nostra, per inceptos himenaeos. Pellentesque ac diam a enim posuere ultricies a in mauris. Pellentesque suscipit, velit vel luctus ultrices, lacus nunc pellentesque ipsum, sed mollis erat est non nunc. Mauris eget ante congue nibh dapibus porttitor. Mauris lacinia rutrum sodales. Curabitur neque erat, efficitur et dictum quis, semper a magna. Fusce fringilla lectus bibendum, viverra arcu quis, commodo diam. Suspendisse auctor ut diam a rhoncus. Duis dignissim, purus eget hendrerit ultricies, sem urna ultricies est, ut maximus tellus felis at elit. Sed vestibulum mi vitae massa blandit rutrum. Sed mattis lacus at odio mattis, nec euismod diam dictum. Mauris accumsan diam sed neque commodo, nec ultrices est pulvinar. Cras vulputate mi vel neque eleifend congue. Sed viverra varius dolor, vel eleifend est tristique sit amet. Duis sagittis nibh id tincidunt condimentum. Fusce at eleifend nibh.
Mauris sed lacus et magna dictum dignissim sed eget nibh. Suspendisse potenti. Nunc scelerisque condimentum magna, at facilisis nulla ultrices quis. Nulla in ex id lorem semper luctus. Suspendisse potenti. Interdum et malesuada fames ac ante ipsum primis in faucibus. Interdum et malesuada fames ac ante ipsum primis in faucibus. Quisque arcu ligula, volutpat sed ligula vehicula, posuere dictum quam.
Fusce fermentum massa convallis sapien aliquet fermentum. Maecenas ac ante at orci imperdiet imperdiet nec eget nulla. Curabitur lacinia tincidunt ipsum, id consequat orci. Nullam auctor nec risus at aliquam. Proin non euismod orci, a viverra quam. Mauris mollis placerat commodo. Aliquam posuere nibh arcu, eu volutpat felis finibus ut. Integer nec auctor sem. Donec blandit enim a rutrum porttitor. Mauris sit amet sapien turpis. Morbi rhoncus auctor diam non elementum.
Aenean sed sapien ut eros hendrerit fermentum eget at turpis. Vivamus in nibh vitae odio faucibus porta vitae et ante. Vestibulum ullamcorper lacinia venenatis. Quisque sollicitudin risus at sapien vulputate posuere. Duis elementum vestibulum vestibulum. Phasellus volutpat luctus ullamcorper. Nam egestas nec sapien ut gravida. Etiam porttitor diam ut posuere aliquam.`)
	compressed, err := br.Encode(data)
	if err != nil {
		t.Error(err)
	}

	logrus.Info(len(data), len(compressed), len(data)-len(compressed))
	decompressed, err := br.Decode(compressed)
	if err != nil {
		t.Error(err)
	}
	if string(decompressed) != string(data) {
		t.Error("decompressed data is not equal to original data")
	}
}
