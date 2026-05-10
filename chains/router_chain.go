package chains

import (
	"strings"
)

// RouteDestination represents a routing target with name, description, and chain.
type RouteDestination struct {
	Name        string
	Description string
	Chain       BaseChain
	Keywords    []string
}

func NewRouteDestination(name, description string, chain BaseChain) *RouteDestination {
	return &RouteDestination{
		Name:        name,
		Description: description,
		Chain:       chain,
	}
}

func (d *RouteDestination) WithKeywords(keywords ...string) *RouteDestination {
	d.Keywords = keywords
	return d
}

// RouterChain routes input to different chains based on keyword matching.
type RouterChain struct {
	destinations []*RouteDestination
	defaultChain BaseChain
	InputKey     string
	name         string
	Verbose      bool
}

func NewRouterChain() *RouterChain {
	return &RouterChain{
		destinations: make([]*RouteDestination, 0),
		InputKey:     "input",
		name:         "router_chain",
	}
}

func (c *RouterChain) AddRoute(name, description string, chain BaseChain) *RouterChain {
	c.destinations = append(c.destinations, NewRouteDestination(name, description, chain))
	return c
}

func (c *RouterChain) AddRouteWithKeywords(name, description string, chain BaseChain, keywords []string) *RouterChain {
	c.destinations = append(c.destinations, NewRouteDestination(name, description, chain).WithKeywords(keywords...))
	return c
}

func (c *RouterChain) WithDefault(chain BaseChain) *RouterChain {
	c.defaultChain = chain
	return c
}

func (c *RouterChain) WithInputKey(key string) *RouterChain {
	c.InputKey = key
	return c
}

func (c *RouterChain) WithName(name string) *RouterChain {
	c.name = name
	return c
}

func (c *RouterChain) WithVerbose(v bool) *RouterChain {
	c.Verbose = v
	return c
}

func (c *RouterChain) InputKeys() []string {
	return []string{c.InputKey}
}

func (c *RouterChain) OutputKeys() []string {
	if c.defaultChain != nil {
		return c.defaultChain.OutputKeys()
	}
	if len(c.destinations) > 0 {
		return c.destinations[0].Chain.OutputKeys()
	}
	return []string{"output"}
}

func (c *RouterChain) ValidateInputs(inputs map[string]any) error {
	if _, ok := inputs[c.InputKey]; !ok {
		return NewChainError(ErrMissingInput, "missing input: "+c.InputKey, nil)
	}
	return nil
}

func (c *RouterChain) routeByKeywords(input string) *RouteDestination {
	for _, dest := range c.destinations {
		for _, kw := range dest.Keywords {
			if strings.Contains(input, kw) {
				return dest
			}
		}
	}
	return nil
}

func (c *RouterChain) selectRoute(input string) (*RouteDestination, error) {
	if dest := c.routeByKeywords(input); dest != nil {
		return dest, nil
	}
	return nil, nil
}

func (c *RouterChain) Invoke(inputs map[string]any) (ChainResult, error) {
	if err := c.ValidateInputs(inputs); err != nil {
		return nil, err
	}

	input, _ := inputs[c.InputKey].(string)
	dest, _ := c.selectRoute(input)

	if dest != nil {
		return dest.Chain.Invoke(inputs)
	}

	if c.defaultChain != nil {
		return c.defaultChain.Invoke(inputs)
	}

	return nil, NewChainError(ErrExecution,
		"no matching route and no default chain configured", nil)
}

func (c *RouterChain) Name() string {
	return c.name
}

// LLMRouterChain uses a RouterChain for keyword-based routing.
// LLM-based routing requires external LLM integration.
type LLMRouterChain struct {
	inner *RouterChain
}

func NewLLMRouterChain() *LLMRouterChain {
	return &LLMRouterChain{inner: NewRouterChain()}
}

func (c *LLMRouterChain) AddRoute(name, description string, chain BaseChain) *LLMRouterChain {
	c.inner.AddRoute(name, description, chain)
	return c
}

func (c *LLMRouterChain) AddRouteWithKeywords(name, description string, chain BaseChain, keywords []string) *LLMRouterChain {
	c.inner.AddRouteWithKeywords(name, description, chain, keywords)
	return c
}

func (c *LLMRouterChain) WithDefault(chain BaseChain) *LLMRouterChain {
	c.inner.WithDefault(chain)
	return c
}

func (c *LLMRouterChain) WithInputKey(key string) *LLMRouterChain {
	c.inner.WithInputKey(key)
	return c
}

func (c *LLMRouterChain) InputKeys() []string  { return c.inner.InputKeys() }
func (c *LLMRouterChain) OutputKeys() []string { return c.inner.OutputKeys() }
func (c *LLMRouterChain) ValidateInputs(inputs map[string]any) error { return c.inner.ValidateInputs(inputs) }
func (c *LLMRouterChain) Invoke(inputs map[string]any) (ChainResult, error) { return c.inner.Invoke(inputs) }
func (c *LLMRouterChain) Name() string { return "llm_router_chain" }
