const TOTAL = 203
const LENGTH = 50
for (let idx = 0; idx < TOTAL; idx++) {
  const array: number[] = []

  let jdx = idx
  for (; jdx < idx + LENGTH; jdx++) {
    if (jdx > TOTAL) {
      break
    }
    array.push(jdx)
  }

  console.log(array)
  idx = jdx - 1
}
