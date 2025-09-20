package packsizes

type Provider interface {
	List() ([]int, error)
}
