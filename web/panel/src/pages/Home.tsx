import { Box, Center, Space } from "@mantine/core";
import { HomeIcon } from "lucide-react";

export default function Home() {
  return (
    <Center mih={500}>
      <HomeIcon />
      <Space w="xs" />
      <Box>Home</Box>
    </Center>
  )
}
