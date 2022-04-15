package priorities

/*
type resource struct {
	pod   *v1.Pod
	nodes []*v1.Node
	name  string
}

func TestNodeAffinityPriority(t *testing.T) {
	label1 := map[string]string{"foo": "bar"}
	label2 := map[string]string{"key": "value"}
	label3 := map[string]string{"az": "az1"}
	// label4 := map[string]string{"abc": "az11", "def": "az22"}
	// label5 := map[string]string{"foo": "bar", "key": "value", "az": "az1"}

	affinity1 := &v1.Affinity{
		NodeAffinity: &v1.NodeAffinity{
			PreferredDuringSchedulingIgnoredDuringExecution: []v1.PreferredSchedulingTerm{{
				Weight: 2,
				Preference: v1.NodeSelectorTerm{
					MatchExpressions: []v1.NodeSelectorRequirement{{
						Key:      "foo",
						Operator: v1.NodeSelectorOpIn,
						Values:   []string{"bar"},
					}},
				},
			}},
		},
	}
	_ = affinity1
	affinity2 := &v1.Affinity{
		NodeAffinity: &v1.NodeAffinity{
			PreferredDuringSchedulingIgnoredDuringExecution: []v1.PreferredSchedulingTerm{
				{
					Weight: 2,
					Preference: v1.NodeSelectorTerm{
						MatchExpressions: []v1.NodeSelectorRequirement{
							{
								Key:      "foo",
								Operator: v1.NodeSelectorOpIn,
								Values:   []string{"bar"},
							},
						},
					},
				},
				{
					Weight: 4,
					Preference: v1.NodeSelectorTerm{
						MatchExpressions: []v1.NodeSelectorRequirement{
							{
								Key:      "key",
								Operator: v1.NodeSelectorOpIn,
								Values:   []string{"value"},
							},
						},
					},
				},
				{
					Weight: 5,
					Preference: v1.NodeSelectorTerm{
						MatchExpressions: []v1.NodeSelectorRequirement{
							{
								Key:      "foo",
								Operator: v1.NodeSelectorOpIn,
								Values:   []string{"bar"},
							},
							{
								Key:      "key",
								Operator: v1.NodeSelectorOpIn,
								Values:   []string{"value"},
							},
							{
								Key:      "az",
								Operator: v1.NodeSelectorOpIn,
								Values:   []string{"az1"},
							},
						},
					},
				},
			},
		},
	}

	test := resource{
		pod: &v1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Annotations: map[string]string{},
			},
			Spec: v1.PodSpec{
				Affinity: affinity2,
			},
		},
		nodes: []*v1.Node{
			{ObjectMeta: metav1.ObjectMeta{Name: "machine1", Labels: label1}},
			{ObjectMeta: metav1.ObjectMeta{Name: "machine2", Labels: label2}},
			{ObjectMeta: metav1.ObjectMeta{Name: "machine3", Labels: label3}},
		},

		name: "all machines are same priority as NodeAffinity is nil",
	}

	t.Run(test.name, func(t *testing.T) {
		for i, r := range test.nodes {
			score := NodeAffinity(test.pod, r)
			fmt.Println(i, score)
		}
	})

}
*/
