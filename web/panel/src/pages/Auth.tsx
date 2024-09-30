import { Badge, Box, Button, Center, Group, PasswordInput, rem, Stack, Text } from "@mantine/core";
import { FileKey2, Upload, X, ChevronRight, LogIn } from "lucide-react";
import { Dropzone } from '@mantine/dropzone';

const iconSz = { width: rem(52), height: rem(52) };

export default function Auth() {
  return (
    <Center
      p="md"
    >
      <Stack
        py="xl"
        gap="xl"
        w={"90%"}
        align="center"
      >
        <Badge
          size="lg"
          tt="none"
          variant="gradient"
          gradient={{ from: 'teal', to: 'green', deg: 209 }}
        >
          D2M is running on your machine
        </Badge>

        <Dropzone
          w="100%"
          maw={600}
          accept={{ "application/x-x509-ca-cert": [".pem"] }}
          onDrop={(files) => {
            console.log(files);
          }}
        >
          <Group justify="center" gap="xl" mih={150} style={{ pointerEvents: 'none' }}>
            <Dropzone.Accept>
              <Upload
                style={{ ...iconSz, color: 'var(--mantine-color-blue-6)' }}
              />
            </Dropzone.Accept>
            <Dropzone.Reject>
              <X
                style={{  ...iconSz, color: 'var(--mantine-color-red-6)' }}
              />
            </Dropzone.Reject>
            <Dropzone.Idle>
              <FileKey2
                style={{  ...iconSz, color: 'var(--mantine-color-dimmed)' }}
              />
            </Dropzone.Idle>

            <Box>
              <Text size="lg" lh={1.2}>
                Drag n Drop or click here to upload you public key file
              </Text>
              <Text size="sm" c="dimmed" lh={1.2} mt={7}>
                (.pem File) Configured in D2M CLI
              </Text>
            </Box>
          </Group>
        </Dropzone>
        <PasswordInput
          w="100%"
          maw={600}
          size="md"
          variant="filled"
          label="Access Password?"
          placeholder="Enter it here...."
          description="You have configured it for your D2M instance in the CLI."
        />
        <Button
          leftSection={<LogIn size={16} />}
          rightSection={<ChevronRight />}
          loaderProps={{ type: "bars" }}
          variant="light"
        >
          Authorize
        </Button>
      </Stack>
      
    </Center>
  )
}
