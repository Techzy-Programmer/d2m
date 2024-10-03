import { Center, Loader } from '@mantine/core'

export default function Loading() {
  return (
    <>
      <Center mih="calc(100vh - 100px)">
        <Loader type='bars' />
      </Center>
    </>
  )
}
