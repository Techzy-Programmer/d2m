import { Badge, Box, Button, Center, Group, PasswordInput, rem, Stack, Text } from "@mantine/core";
import { FileKey2, Upload, X, ChevronRight, LogIn } from "lucide-react";
import { isFetchSuccess, showToast } from "../utils/general";
import { Dropzone, FileWithPath } from '@mantine/dropzone';
import useEncryption from "../hooks/useEncryption";
import { useNavigate } from "react-router-dom";
import { useAuth } from "../state/use-auth";
import { useEffect, useState } from "react";
import { MetaResp } from "../utils/types";
import useFetch from "../hooks/useFetch";
import { useMeta } from "../state/use-meta";

const iconSz = { width: rem(52), height: rem(52) };

export default function Auth() {
  const fetchData = useFetch();
  const navigate = useNavigate();
  const { setLoggedIn } = useAuth();
  const [busy, setBusy] = useState(false);
  const [pubKey, setPubKey] = useState('');
  const [accessPwd, setAccessPwd] = useState('');
  const { encError, encrypt } = useEncryption(pubKey);
  const { setPageTitle, setMetadata } = useMeta();
  const hasRSAError = encError;

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

    const payload = encrypt(accessPwd);

    if (payload === null) {
      setBusy(false);
      return showToast({
        title: "Encryption Failed",
        message: "Internal encryption error",
        status: "issue",
      });
    }

    const fetch = await fetchData<MetaResp>("/api/auth", {
      headers: { "X-Encryption-Key": payload.encKey },
      body: payload.encData,
      method: "POST",
    });

    setBusy(false);

    if (!isFetchSuccess(fetch)) {
      return showToast({
        message: fetch.error || "Something went wrong",
        title: "Failed To Fetch",
        status: "issue",
      });
    }

    const { code, data } = fetch;

    if (code !== 200) {
      return showToast({
        title: "Authentication Failed",
        message: data.message,
        status: "warn",
      });
    }

    showToast({
      title: "Authentication Successful",
      message: data.message,
      status: "ok",
    });

    setMetadata(data.meta);
    setPageTitle("Home");
    setLoggedIn(true);
    navigate("/");
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
