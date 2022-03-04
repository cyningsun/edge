package hash

import (
	"fmt"
	"math"
	"reflect"
	"sort"
	"sync"
	"testing"
)

var testRingData = map[string]*ring{
	"zeroNodeRing": {
		replicas: 2,
		sorted:   []uint32{},
		vnodes:   map[uint32]Node{},
		mtx:      sync.Mutex{},
	},
	"oneNodeRing": {
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
	"fourNodeRing": {
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
	"fiveNodeRing": {
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
	"testNode1": {val: "testNode1"},
	"testNode2": {val: "testNode2"},
	"testNode3": {val: "testNode3"},
	"testNode4": {val: "testNode4"},
	"testNode5": {val: "testNode5"},
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
			got, err := NewRing(tt.args.replicas)
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
				mtx:      sync.Mutex{},
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
				mtx:      sync.Mutex{},
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
				mtx:      sync.Mutex{},
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
			ring1, _ := NewRing(tt.args.replicas)
			ring2, _ := NewRing(tt.args.replicas)
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
	hash, _ := NewRing(50)

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

// TestKeyDistribution ref:https://arxiv.org/pdf/1406.2294.pdf
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
			replicas: 10,
			nodeCnt:  3,
		},
		{
			replicas: 100,
			nodeCnt:  3,
		},
		{
			replicas: 1000,
			nodeCnt:  3,
		},
		{
			replicas: 1,
			nodeCnt:  10,
		},
		{
			replicas: 10,
			nodeCnt:  10,
		},
		{
			replicas: 100,
			nodeCnt:  10,
		},
		{
			replicas: 1000,
			nodeCnt:  10,
		},
		{
			replicas: 1,
			nodeCnt:  20,
		},
		{
			replicas: 10,
			nodeCnt:  20,
		},
		{
			replicas: 100,
			nodeCnt:  20,
		},
		{
			replicas: 1000,
			nodeCnt:  20,
		},
		{
			replicas: 1,
			nodeCnt:  50,
		},
		{
			replicas: 10,
			nodeCnt:  50,
		},
		{
			replicas: 100,
			nodeCnt:  50,
		},
		{
			replicas: 1000,
			nodeCnt:  50,
		},
		{
			replicas: 1,
			nodeCnt:  100,
		},
		{
			replicas: 10,
			nodeCnt:  100,
		},
		{
			replicas: 100,
			nodeCnt:  100,
		},
		{
			replicas: 1000,
			nodeCnt:  100,
		},
	}

	fmt.Print("| Name | Standard Error | 99% Confidence Interval | \n")
	fmt.Print("|----|----|\n")
	for _, tt := range tests {
		h, _ := NewRing(tt.replicas)
		for i := 0; i < tt.nodeCnt; i++ {
			h.Add(testNode{val: fmt.Sprintf("node-%d", i)})
		}

		totalBuckets := tt.replicas * tt.nodeCnt
		allBucketInterval := make([]uint32, 0, totalBuckets)

		bucketInterval := math.MaxUint32 - h.sorted[totalBuckets-1] + h.sorted[0]
		allBucketInterval = append(allBucketInterval, bucketInterval)
		for i := 1; i < totalBuckets; i++ {
			bucketInterval = h.sorted[i] - h.sorted[i-1]
			allBucketInterval = append(allBucketInterval, bucketInterval)
		}

		stdDev := StandardDeviation(allBucketInterval)
		logBucket := math.Sqrt(float64(totalBuckets))
		standErr := stdDev / logBucket

		idealInterval := Mean(allBucketInterval)
		lower, upper := NormalConfidenceInterval(allBucketInterval)
		fmt.Printf("|%v(nodeCnt)- %v(replicas)| %7f  | (%7f,%7f)|\n", tt.nodeCnt, tt.replicas, standErr/idealInterval, lower/idealInterval, upper/idealInterval)
	}
}

// StandardDeviation returns the standard deviation of the slice
// as a float
func StandardDeviation(nums []uint32) (dev float64) {
	if len(nums) == 0 {
		return 0.0
	}

	m := Mean(nums)
	for _, n := range nums {
		dev += (float64(n) - m) * (float64(n) - m)
	}
	dev = math.Pow(dev/float64(len(nums)), 0.5)
	return dev
}

// NormalConfidenceInterval returns the 99% confidence interval for the mean
// as two float values, the lower and the upper bounds and assuming a normal
// distribution
func NormalConfidenceInterval(nums []uint32) (lower, upper float64) {
	conf := 2.57583 // 99% confidence for the mean, http://bit.ly/Mm05eZ
	mean := Mean(nums)
	dev := StandardDeviation(nums) / math.Sqrt(float64(len(nums)))
	return mean - dev*conf, mean + dev*conf
}

// Mean returns the mean of an integer array as a float
func Mean(nums []uint32) (mean float64) {
	if len(nums) == 0 {
		return 0.0
	}
	for _, n := range nums {
		mean += float64(n)
	}
	return mean / float64(len(nums))
}
