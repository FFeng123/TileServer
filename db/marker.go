package db

type Marker [3]string

func (m Marker) String() string {
	return m[0] + "," + m[1] + "," + m[2]
}

func NewMarker(x, y, z string) Marker {
	return Marker{x, y, z}
}
