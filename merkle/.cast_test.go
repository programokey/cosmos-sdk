package merkle

import (
	"fmt"
	"math/rand"
	"testing"

	dbm "github.com/tendermint/tmlibs/db"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tendermint/iavl"
)

func TestExistsProof(t *testing.T) {
	db := dbm.NewMemDB()
	tree := iavl.NewVersionedTree(db, 0)

	size := 2

	keys := make([][]byte, size)
	vals := make([][]byte, size)

	for i := 0; i < size; i++ {
		keys[i] = make([]byte, rand.Uint32()%64)
		rand.Read(keys[i])
		vals[i] = make([]byte, rand.Uint32()%64)
		rand.Read(vals[i])

		tree.Set(keys[i], vals[i])
	}

	root, ver, err := tree.SaveVersion()
	require.Nil(t, err)

	fmt.Printf("%+v, %+v\n", root, ver)

	for i, k := range keys {
		fmt.Printf("\n\n#%d\n", i)
		val, prf, err := tree.GetVersionedWithProof(k, ver)
		require.Nil(t, err)

		require.Nil(t, prf.Verify(k, vals[i], root))

		assert.Equal(t, vals[i], val)

		mprf, err := FromKeyProof(prf)
		assert.Nil(t, err)

		fmt.Printf("%+v, %+v\n", k, vals[i])
		leaf, err := Leaf(k, vals[i])
		assert.Nil(t, err)

		fmt.Printf("%+v\n", leaf)
		calcroot, err := mprf.Run(leaf)
		assert.Nil(t, err, calcroot)
		fmt.Printf("%+v\n", prf)
		fmt.Printf("%+v\n", mprf)
		assert.Equal(t, root, calcroot)
	}
}
