package ringhash

import (
	"fmt"
	"reflect"
	"sort"
	"sync"
	"testing"
)

var testRingData = map[string]*ring{
	"zeroNodeRing": &ring{
		replicas: 2,
		sorted:   []uint32{},
		vnodes:   map[uint32]Node{},
		mtx:      sync.Mutex{},
	},
	"oneNodeRing": &ring{
		replicas: 2,
		sorted: sorted([]uint32{
			hash("testNode1_1"),
			hash("testNode1_2"),
		}),
		vnodes: map[uint32]Node{
			hash("testNode1_1"): testNodeData["testNode1"],
			hash("testNode1_2"): testNodeData["testNode1"],
		},
		mtx: sync.Mutex{},
	},
	"fourNodeRing": &ring{
		replicas: 2,
		sorted: sorted([]uint32{
			hash("testNode1_1"),
			hash("testNode2_1"),
			hash("testNode3_1"),
			hash("testNode4_1"),

			hash("testNode1_2"),
			hash("testNode2_2"),
			hash("testNode3_2"),
			hash("testNode4_2"),
		}),
		vnodes: map[uint32]Node{
			hash("testNode1_1"): testNodeData["testNode1"],
			hash("testNode2_1"): testNodeData["testNode2"],
			hash("testNode3_1"): testNodeData["testNode3"],
			hash("testNode4_1"): testNodeData["testNode4"],

			hash("testNode1_2"): testNodeData["testNode1"],
			hash("testNode2_2"): testNodeData["testNode2"],
			hash("testNode3_2"): testNodeData["testNode3"],
			hash("testNode4_2"): testNodeData["testNode4"],
		},
		mtx: sync.Mutex{},
	},
	"fiveNodeRing": &ring{
		replicas: 2,
		sorted: sorted([]uint32{
			hash("testNode1_1"),
			hash("testNode2_1"),
			hash("testNode3_1"),
			hash("testNode4_1"),
			hash("testNode5_1"),

			hash("testNode1_2"),
			hash("testNode2_2"),
			hash("testNode3_2"),
			hash("testNode4_2"),
			hash("testNode5_2"),
		}),
		vnodes: map[uint32]Node{
			hash("testNode1_1"): testNodeData["testNode1"],
			hash("testNode2_1"): testNodeData["testNode2"],
			hash("testNode3_1"): testNodeData["testNode3"],
			hash("testNode4_1"): testNodeData["testNode4"],
			hash("testNode5_1"): testNodeData["testNode5"],

			hash("testNode1_2"): testNodeData["testNode1"],
			hash("testNode2_2"): testNodeData["testNode2"],
			hash("testNode3_2"): testNodeData["testNode3"],
			hash("testNode4_2"): testNodeData["testNode4"],
			hash("testNode5_2"): testNodeData["testNode5"],
		},
		mtx: sync.Mutex{},
	},
}

var testNodeData = map[string]testNode{
	"testNode1": testNode{val: "testNode1"},
	"testNode2": testNode{val: "testNode2"},
	"testNode3": testNode{val: "testNode3"},
	"testNode4": testNode{val: "testNode4"},
	"testNode5": testNode{val: "testNode5"},
}

func sorted(slice []uint32) []uint32 {
	sort.Slice(slice, func(i, j int) bool { return slice[i] < slice[j] })
	return slice
}

func TestNew(t *testing.T) {
	type args struct {
		replicas int
	}
	tests := []struct {
		name    string
		args    args
		want    *ring
		wantErr bool
	}{
		{
			"normal",
			args{
				replicas: 2,
			},
			testRingData["zeroNodeRing"],
			false,
		},
		{
			"min replicas validation",
			args{
				replicas: 0,
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.replicas)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

type testNode struct {
	val string
}

func (t testNode) String() string {
	return t.val
}

func Test_ring_Add(t *testing.T) {
	type fields struct {
		replicas int
		nodes    map[uint32]Node
		sorted   []uint32
		mtx      sync.Mutex
	}
	type args struct {
		node Node
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *ring
	}{
		{
			"normal",
			fields{
				replicas: testRingData["zeroNodeRing"].replicas,
				nodes:    testRingData["zeroNodeRing"].vnodes,
				sorted:   testRingData["zeroNodeRing"].sorted,
				mtx:      testRingData["zeroNodeRing"].mtx,
			},
			args{
				node: testNodeData["testNode1"],
			},
			testRingData["oneNodeRing"],
		},
		{
			"duplicate add",
			fields{
				replicas: testRingData["oneNodeRing"].replicas,
				nodes:    testRingData["oneNodeRing"].vnodes,
				sorted:   testRingData["oneNodeRing"].sorted,
				mtx:      testRingData["oneNodeRing"].mtx,
			},
			args{
				node: testNodeData["testNode1"],
			},
			testRingData["oneNodeRing"],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := &ring{
				replicas: tt.fields.replicas,
				vnodes:   tt.fields.nodes,
				sorted:   tt.fields.sorted,
				mtx:      tt.fields.mtx,
			}
			got.Add(tt.args.node)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_hash(t *testing.T) {
	type args struct {
		n string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"normal",
			args{
				"node_1",
			},
			"e1404e65", // [fnv.32a("node_1")](https://md5calc.com/hash/fnv1a32?str=node_1)
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fmt.Sprintf("%x", hash(tt.args.n)); got != tt.want {
				t.Errorf("hash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ring_Get(t *testing.T) {
	type fields struct {
		replicas int
		nodes    map[uint32]Node
		sorted   []uint32
		mtx      sync.Mutex
	}
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Node
	}{
		{
			"normal",
			fields{
				replicas: testRingData["fiveNodeRing"].replicas,
				nodes:    testRingData["fiveNodeRing"].vnodes,
				sorted:   testRingData["fiveNodeRing"].sorted,
				mtx:      testRingData["fiveNodeRing"].mtx,
			},
			args{
				"key",
			},
			testNodeData["testNode2"],
		},
		{
			"zore node",
			fields{
				replicas: testRingData["zeroNodeRing"].replicas,
				nodes:    testRingData["zeroNodeRing"].vnodes,
				sorted:   testRingData["zeroNodeRing"].sorted,
				mtx:      testRingData["zeroNodeRing"].mtx,
			},
			args{
				"key",
			},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &ring{
				replicas: tt.fields.replicas,
				vnodes:   tt.fields.nodes,
				sorted:   tt.fields.sorted,
				mtx:      tt.fields.mtx,
			}
			if got := r.Get(tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ring.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ring_Remove(t *testing.T) {
	type fields struct {
		replicas int
		nodes    map[uint32]Node
		sorted   []uint32
		mtx      sync.Mutex
	}
	type args struct {
		node Node
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *ring
	}{
		{
			"normal",
			fields{
				replicas: testRingData["fiveNodeRing"].replicas,
				nodes:    testRingData["fiveNodeRing"].vnodes,
				sorted:   testRingData["fiveNodeRing"].sorted,
				mtx:      testRingData["fiveNodeRing"].mtx,
			},
			args{
				testNodeData["testNode5"],
			},
			testRingData["fourNodeRing"],
		},
		{
			"remove not exist",
			fields{
				replicas: testRingData["oneNodeRing"].replicas,
				nodes:    testRingData["oneNodeRing"].vnodes,
				sorted:   testRingData["oneNodeRing"].sorted,
				mtx:      testRingData["oneNodeRing"].mtx,
			},
			args{
				testNodeData["testNode2"],
			},
			testRingData["oneNodeRing"],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := &ring{
				replicas: tt.fields.replicas,
				vnodes:   tt.fields.nodes,
				sorted:   tt.fields.sorted,
				mtx:      tt.fields.mtx,
			}
			got.Remove(tt.args.node)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConsistency(t *testing.T) {
	type args struct {
		replicas int
		nodes1   []Node
		nodes2   []Node
		key      string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"random key",
			args{
				2,
				[]Node{
					testNodeData["testNode1"],
					testNodeData["testNode2"],
				},
				[]Node{
					testNodeData["testNode2"],
					testNodeData["testNode1"],
				},
				"testNode",
			},
			true,
		},
		{
			"vnode1 key",
			args{
				1,
				[]Node{
					testNodeData["testNode1"],
					testNodeData["testNode2"],
				},
				[]Node{
					testNodeData["testNode2"],
					testNodeData["testNode1"],
				},
				"testNode1_1",
			},
			true,
		},
		{
			"vnode2 key",
			args{
				2,
				[]Node{
					testNodeData["testNode1"],
					testNodeData["testNode2"],
				},
				[]Node{
					testNodeData["testNode2"],
					testNodeData["testNode1"],
				},
				"testNode1_2",
			},
			true,
		},
		{
			"vnode3 key",
			args{
				2,
				[]Node{
					testNodeData["testNode1"],
					testNodeData["testNode2"],
				},
				[]Node{
					testNodeData["testNode2"],
					testNodeData["testNode1"],
				},
				"testNode2_1",
			},
			true,
		},
		{
			"vnode4 key",
			args{
				2,
				[]Node{
					testNodeData["testNode1"],
					testNodeData["testNode2"],
				},
				[]Node{
					testNodeData["testNode2"],
					testNodeData["testNode1"],
				},
				"testNode2_2",
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ring1, _ := New(tt.args.replicas)
			ring2, _ := New(tt.args.replicas)
			for _, each := range tt.args.nodes1 {
				ring1.Add(each)
			}
			for _, each := range tt.args.nodes2 {
				ring2.Add(each)
			}

			node1 := ring1.Get(tt.args.key)
			node2 := ring2.Get(tt.args.key)
			if reflect.DeepEqual(node1, node2) != tt.want {
				t.Errorf("Rings consistent check failed, key:%v, node1:%v, node2:%v, want:%v", tt.args.key, node1, node2, tt.want)
			}
		})
	}
}

func BenchmarkGet8(b *testing.B)   { benchmarkGet(b, 8) }
func BenchmarkGet32(b *testing.B)  { benchmarkGet(b, 32) }
func BenchmarkGet128(b *testing.B) { benchmarkGet(b, 128) }
func BenchmarkGet512(b *testing.B) { benchmarkGet(b, 512) }

func benchmarkGet(b *testing.B, nodes int) {

	hash, _ := New(50)

	var buckets []testNode
	for i := 0; i < nodes; i++ {
		buckets = append(buckets, testNode{val: fmt.Sprintf("node-%d", i)})
	}

	for _, each := range buckets {
		hash.Add(each)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		hash.Get(buckets[i&(nodes-1)].val)
	}
}

// https://arxiv.org/pdf/1406.2294.pdf
func TestKeyDistribution(t *testing.T) {
	tests := []struct {
		replicas int
		nodeCnt  int
	}{
		{
			replicas: 1,
			nodeCnt:  3,
		},
		{
			replicas: 1,
			nodeCnt:  10,
		},
		{
			replicas: 1,
			nodeCnt:  20,
		},
		{
			replicas: 1,
			nodeCnt:  50,
		},
		{
			replicas: 1,
			nodeCnt:  100,
		},
		{
			replicas: 10,
			nodeCnt:  3,
		},
		{
			replicas: 10,
			nodeCnt:  10,
		},
		{
			replicas: 10,
			nodeCnt:  20,
		},
		{
			replicas: 10,
			nodeCnt:  50,
		},
		{
			replicas: 10,
			nodeCnt:  100,
		},
		{
			replicas: 100,
			nodeCnt:  3,
		},
		{
			replicas: 100,
			nodeCnt:  10,
		},
		{
			replicas: 100,
			nodeCnt:  20,
		},
		{
			replicas: 100,
			nodeCnt:  50,
		},
		{
			replicas: 100,
			nodeCnt:  100,
		},
		{
			replicas: 1000,
			nodeCnt:  3,
		},
		{
			replicas: 1000,
			nodeCnt:  10,
		},
		{
			replicas: 1000,
			nodeCnt:  20,
		},
		{
			replicas: 1000,
			nodeCnt:  50,
		},
		{
			replicas: 1000,
			nodeCnt:  100,
		},
	}

	for _, tt := range tests {
		h, _ := New(tt.replicas)
		for i := 0; i < tt.nodeCnt; i++ {
			h.Add(testNode{val: fmt.Sprintf("node-%d", i)})
		}

		// totalBuckets := tt.replicas * tt.nodeCnt
		// idealInterval := math.MaxUint32 / uint32(totalBuckets)
		// intervalDeviation := make([]float64, 0, totalBuckets)
		// bucketInterval := math.MaxUint32 - h.sorted[totalBuckets-1] + h.sorted[0]
		// intervalDeviation = append(intervalDeviation, math.Abs(float64(bucketInterval-idealInterval)/float64(idealInterval)))
		// for i := 1; i < totalBuckets; i++ {
		// 	bucketInterval = h.sorted[i] - h.sorted[i-1]
		// 	intervalDeviation = append(intervalDeviation, math.Abs(float64(bucketInterval-idealInterval)/float64(idealInterval)))
		// }
		// if len(intervalDeviation) != totalBuckets {
		// 	t.Fatalf("totalBuckets:%v, bucketInterval:%v, ", totalBuckets, bucketInterval)
		// }
		// stdDev, _ := stats.StandardDeviation(intervalDeviation)
		// logBucket := math.Log10(float64(totalBuckets))
		// mean, _ := stats.Mean(intervalDeviation)
		// standErr := stdDev / mean
		// t.Log("Name | Standard Error | 99% Confidence Interval")
		// t.Logf("%v(replicas)-%v(nodeCnt) | %v  | (%v,%v)\n", tt.replicas, tt.nodeCnt, standErr, mean-2.576*stdDev/logBucket, mean+2.576*stdDev/logBucket)
	}
}
