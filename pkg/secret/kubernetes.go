package secret


type KubernetesSecret struct {
        Metadata map[string]string
        Data map[string]string
}
