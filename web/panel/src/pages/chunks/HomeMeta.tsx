import { Box, Center, Flex, Paper, Text, Title } from "@mantine/core";
import { MetaResp } from "../../utils/types";

type MetaProps = MetaResp["meta"];

export default function HomeMeta({ uptime, webPort, tcpPort }: MetaProps) {
  return (
    <Center mt={10}>
      <Paper
        withBorder
        w="90%"
        p="md"
        m="sm"
      >
        <Flex
          justify="start"
          align="start"
          wrap="wrap"
          gap={20}
        >
          <Box>
            <Title
              tt="capitalize"
              order={5}
              c="gray"
              lh={1.2}
            >
              Ports Occupied?
            </Title>
            <Text
              gradient={{ from: 'cyan', to: 'lime', deg: 180 }}
              variant="gradient"
              ff="monospace"
              size="lg"
              fw={900}
            >
              {webPort} & {tcpPort}
            </Text>
          </Box>
          <Box>
            <Title
              tt="capitalize"
              order={5}
              c="gray"
              lh={1.2}
            >
              Uptime?
            </Title>
            <Text
              gradient={{ from: 'cyan', to: 'lime', deg: 180 }}
              variant="gradient"
              ff="monospace"
              size="lg"
              fw={900}
            >
              {uptime}
            </Text>
          </Box>
        </Flex>
      </Paper>
    </Center>
  )
}
