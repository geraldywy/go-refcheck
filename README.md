https://pdos.csail.mit.edu/6.824/papers/tour-faq.txt

Q: I'm still a little confused about when to choose value or pointer
receivers. Can you provide any concrete/real-world examples of when we
would choose one over the other?

A: When you want to modify the state of the receiver, you have to use
pointer receivers.  If the struct is very big, you probably want to
use a pointer receiver because value receivers operate on a copy.  If
neither applies, you can use a value receiver.  However, be careful
with value receivers; e.g., if you have a mutex in a struct, you
cannot make it a value receiver, because the mutex would be copied,
defeating its purpose.