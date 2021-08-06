package main

// func GenerateSetters(config Config) {
//     s := ""
//     if singular {
//         s += fmt.Sprintf(`
// func (%s *%s) Set%s(%s %s) *%s {
//     %s.%s = []%s{%s}
//     return %s
// }
// `, param, structure.name, singularName, LowerInitial(singularName), field.ty[2:], structure.name, param, field.name, field.ty[2:], LowerInitial(singularName), param)
//     } else {
//         s += fmt.Sprintf(`
// func (%s *%s) Set%s(%s %s) *%s {
//     %s.%s = %s
//     return %s
// }
// `, param, structure.name, title, field.name, field.ty, structure.name, param, field.name, field.name, param)
//     }
// }
