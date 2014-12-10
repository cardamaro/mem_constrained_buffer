mem_constrained_buffer
======================

A simple buffer that is optionally backed by a concrete file on disk when it grows too large.

    import (
      "github.com/cardamaro/mem_constrained_buffer"
      "strings"
    )
    
    func main() {
      r := strings.NewReader("foobar")
      
      buf := mem_constrained_buffer.NewDefaultMemoryConstrainedBuffer()
      defer buf.Close()  // removes any temp files
      
      if _, err := buf.ReadFrom(r); err != nil {
        panic(err)
      }
      
      fmt.Printf("Read %d bytes", buf.Len())
    }
