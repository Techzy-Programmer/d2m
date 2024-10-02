import { Badge, Box, Button, Center, Group, PasswordInput, rem, Stack, Text } from "@mantine/core";
import { FileKey2, Upload, X, ChevronRight, LogIn } from "lucide-react";
import { Dropzone, FileWithPath } from '@mantine/dropzone';
import useRSA from "../hooks/useRSAEncryption";
import { useEffect, useState } from "react";
import useFetch from "../hooks/useFetch";

const iconSz = { width: rem(52), height: rem(52) };

export default function Auth() {
  const fetchData = useFetch();
  const [busy, setBusy] = useState(false);
  const [pubKey, setPubKey] = useState('');
  const [accessPwd, setAccessPwd] = useState('');
  const { rsaEncError, encrypt } = useRSA(pubKey);
  const hasRSAError = rsaEncError;

  useEffect(() => {
    if (pubKey === "") return;
    encrypt("test");
  }, [pubKey]);

  async function handleDrop([file]: FileWithPath[]) {
    const pubKey = await file.text();
    setPubKey(pubKey);
  }

  async function handleAuth() {
    if (hasRSAError || accessPwd === "") return;
    setBusy(true);

    const { data, error } = await fetchData("/api/auth", {
      body: encrypt(accessPwd),
      method: "POST",
    });

    setBusy(false);
    if (error) {
      console.log(error);
    } else {
      console.log(data);
    }
  }

  return (
    <Center p="md">
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
          maxFiles={1}
          disabled={busy}
          multiple={false}
          maxSize={1024 * 5}
          onDrop={handleDrop}
          accept={{ "application/x-x509-ca-cert": [".pem"] }}
          mod={{ status: pubKey === "" ? null : hasRSAError ? "warn" : "ok" }}
        >
          <Group justify="center" gap="xl" mih={150} style={{ pointerEvents: 'none' }}>
            <Dropzone.Accept>
              <Upload style={{ ...iconSz, color: 'var(--mantine-color-blue-6)' }} />
            </Dropzone.Accept>
            <Dropzone.Reject>
              <X style={{  ...iconSz, color: 'var(--mantine-color-red-6)' }} />
            </Dropzone.Reject>
            <Dropzone.Idle>
              <FileKey2 style={{  ...iconSz, color: 'var(--mantine-color-dimmed)' }} />
            </Dropzone.Idle>

            {pubKey === "" ? (
              <Box>
                <Text size="lg" lh={1.2}>
                  Drag n Drop or click here to upload you public key file
                </Text>
                <Text size="sm" c="dimmed" lh={1.2} mt={7}>
                  (.pem File) Configured in D2M CLI
                </Text>
              </Box>
            ) : (
              !hasRSAError ? (
                <Box>
                  <Text size="lg" lh={1.2}>
                    Public Key Uploaded Successfully
                  </Text>
                  <Text size="sm" c="dimmed" lh={1.2} mt={7}>
                    {pubKey.length} characters
                  </Text>
                </Box>
              ) : (
                <Box>
                  <Text size="lg" lh={1.2}>
                    Invalid File Type
                  </Text>
                  <Text size="sm" c="dimmed" lh={1.2} mt={7}>
                    Only valid .pem files are allowed
                  </Text>
                </Box>
              )
            )}
          </Group>
        </Dropzone>

        <PasswordInput
          w="100%"
          maw={600}
          size="md"
          disabled={busy}
          variant="filled"
          value={accessPwd}
          label="Access Password?"
          placeholder="Enter it here...."
          onChange={(e) => setAccessPwd(e.currentTarget.value)}
          description="You have configured it for your D2M instance in the CLI."
        />

        <Button
          leftSection={<LogIn size={16} />}
          rightSection={<ChevronRight />}
          loaderProps={{ type: "bars" }}
          onClick={handleAuth}
          variant="light"
          loading={busy}
        >
          Authorize
        </Button>
      </Stack>
    </Center>
  )
}
